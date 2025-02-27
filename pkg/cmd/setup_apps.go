package cmd

import (
	"capten/pkg/agent"
	"capten/pkg/app"
	"capten/pkg/cert"
	"capten/pkg/clog"
	"capten/pkg/config"
	"capten/pkg/k8s"

	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "sets up apps cluster for usage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		captenConfig, err := config.GetCaptenConfig()
		if err != nil {
			clog.Logger.Errorf("failed to read capten config, %v", err)
			return
		}

		globalValues, err := app.PrepareGlobalVaules(captenConfig)
		if err != nil {
			clog.Logger.Errorf("applications values preparation failed, %v", err)
			return
		}

		kubeconfigPath := captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.KubeConfigFileName)
		err = k8s.CreateNamespaceIfNotExists(kubeconfigPath, captenConfig.CaptenNamespace)
		if err != nil {
			clog.Logger.Errorf("capten namespace creation failed, %v", err)
			return
		}

		if !captenConfig.SkipAppsDeploy {
			err = app.DeployApps(captenConfig, globalValues, captenConfig.CoreAppGroupsFileName)
			if err != nil {
				clog.Logger.Errorf("%v", err)
				return
			}

		}

		if err := cert.PrepareCerts(captenConfig); err != nil {
			clog.Logger.Errorf("failed to generate certificate, %v", err)
			return
		}
		clog.Logger.Info("Certificates prepared for cluster")

		if err := k8s.CreateOrUpdateCertSecrets(captenConfig); err != nil {
			clog.Logger.Errorf("failed to create secret for certs, %v", err)
			return
		}

		err = k8s.CreateOrUpdateClusterIssuer(captenConfig)
		if err != nil {
			clog.Logger.Errorf("failed to create cstorPoolCluster, %v", err)
			return
		}
		clog.Logger.Info("Configured Certificates on Capten Cluster")

		err = k8s.CreateCStorPoolClusters(captenConfig)
		if err != nil {
			clog.Logger.Errorf("failed to configure storage pool, %v", err)
			return
		}
		clog.Logger.Info("Configured storage pool")

		if !captenConfig.SkipAppsDeploy {
			err = app.DeployApps(captenConfig, globalValues, captenConfig.DefaultAppGroupsFileName)
			if err != nil {
				clog.Logger.Errorf("%v", err)
				return
			}
		}
		clog.Logger.Info("Default Applications Installed")

		clog.Logger.Info("Synchonizing Applications with Cluster Agent")
		if err := agent.SyncInstalledAppConfigsOnAgent(captenConfig); err != nil {
			clog.Logger.Errorf("failed to sync installed apps config in cluster, %v", err)
			return
		}
		clog.Logger.Info("Applications Synchonized with Cluster Agent")

		if captenConfig.StoreCredOnAgent {
			err = agent.StoreCredentials(captenConfig, globalValues)
			if err != nil {
				clog.Logger.Errorf("failed to store cluster credentials, %v", err)
				return
			}
		}

	},
}
