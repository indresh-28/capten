package k8s

import (
	"capten/pkg/config"
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func CreateOrUpdateAgnetCertSecret(captenConfig config.CaptenConfig) error {
	// config, err := clientcmd.BuildConfigFromFlags("", captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName))
	// if err != nil {
	// 	return errors.WithMessage(err, "error while building kubeconfig")
	// }

	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")

	// Build the client configuration
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return errors.WithMessage(err, "Failed to build config")

	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.WithMessage(err, "error while getting k8s config")
	}

	certData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.AgentCertFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading client cert")
	}
	keyData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.AgentKeyFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading client key")
	}
	caCertChainData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.CAFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading ca cert chain")
	}

	// Create the Secret object
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      captenConfig.AgentCertSecretName,
			Namespace: captenConfig.CaptenNamespace,
		},
		Data: map[string][]byte{
			corev1.TLSCertKey:       certData,
			corev1.TLSPrivateKeyKey: keyData,
			"ca.crt":                caCertChainData,
		},
		Type: corev1.SecretTypeTLS,
	}

	_, err = clientSet.CoreV1().Secrets(secret.Namespace).Get(context.Background(), secret.Name, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		_, err := clientSet.CoreV1().Secrets(secret.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
		if err != nil {
			return errors.WithMessage(err, "error in creating secret")
		}
		return nil
	}

	_, err = clientSet.CoreV1().Secrets(secret.Namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
	if err != nil {
		return errors.WithMessage(err, "error in updating secret")
	}
	return nil
}
func CreateOrUpdateCASecret(captenConfig config.CaptenConfig) error {
	config, err := clientcmd.BuildConfigFromFlags("", captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName))
	if err != nil {
		return errors.WithMessage(err, "error while building kubeconfig")
	}
	// home := homedir.HomeDir()
	// kubeconfig := filepath.Join(home, ".kube", "config")

	// // Build the client configuration
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// if err != nil {
	// 	return errors.WithMessage(err, "Failed to build config")

	// }

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.WithMessage(err, "error while getting k8s config")
	}

	interCACertData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.InterCACertFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading client cert")
	}
	interCAKeyData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.InterCAKeyFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading client key")
	}
	caCertChainData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.CAFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading ca cert chain")
	}
	// interCACertData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.InterCACertFileName))
	// if err != nil {
	//     return errors.WithMessage(err, "error while reading inter-ca cert")
	// }
	// interCAKeyData, err := ioutil.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.InterCAKeyFileName))
	// if err != nil {
	//     return errors.WithMessage(err, "error while reading inter-ca key")
	// }

	// Create the Secret object
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      captenConfig.InterCertSecretName,
			Namespace: captenConfig.CaptenNamespace,
		},
		Data: map[string][]byte{
			corev1.TLSCertKey:       interCACertData,
			corev1.TLSPrivateKeyKey: interCAKeyData,
			"ca.crt":                caCertChainData,
			//"inter-ca.crt": interCACertData,
			//"inter-ca.key": interCAKeyData,
		},
		Type: corev1.SecretTypeTLS,
	}

	_, err = clientSet.CoreV1().Secrets(secret.Namespace).Get(context.Background(), secret.Name, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		_, err := clientSet.CoreV1().Secrets(secret.Namespace).Create(context.Background(), secret, metav1.CreateOptions{})
		if err != nil {
			return errors.WithMessage(err, "error in creating secret")
		}
		return nil
	}

	_, err = clientSet.CoreV1().Secrets(secret.Namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
	if err != nil {
		return errors.WithMessage(err, "error in updating secret")
	}
	return nil
}
