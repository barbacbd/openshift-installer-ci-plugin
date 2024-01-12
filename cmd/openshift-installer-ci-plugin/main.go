package openshift_installer_ci_plugin

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/test-infra/prow/config/secret"
	"k8s.io/test-infra/prow/githubeventserver"
	"k8s.io/test-infra/prow/interrupts"
	"k8s.io/test-infra/prow/logrusutil"
	"k8s.io/test-infra/prow/pjutil"
	"os"
	"time"
)

func main() {
	logrusutil.ComponentInit()
	logger := logrus.WithField("plugin", "jira-lifecycle")

	o := gatherOptions()
	if o.validateConfig != "" {
		bytes, err := gzip.ReadFileMaybeGZIP(o.validateConfig)
		if err != nil {
			logger.Fatalf("couldn't read configuration file %s: %v", o.configPath, err)
		}
		if err := validateConfig(bytes); err != nil {
			fmt.Printf("Config is invalid: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if err := o.Validate(); err != nil {
		logger.Fatalf("Invalid options: %v", err)
	}

	configWatchAndUpdate, err := o.getConfigWatchAndUpdate()
	if err != nil {
		logger.WithError(err).Fatal("couldn't get config file watch and update function")
	}
	interrupts.Run(configWatchAndUpdate)

	// get prow config
	configAgent, err := o.prowConfig.ConfigAgent()
	if err != nil {
		logger.WithError(err).Fatal("Error starting config agent.")
	}

	var tokens []string

	// Append the path of hmac and github secrets.
	if o.github.TokenPath != "" {
		tokens = append(tokens, o.github.TokenPath)
	}
	if o.github.AppPrivateKeyPath != "" {
		tokens = append(tokens, o.github.AppPrivateKeyPath)
	}
	tokens = append(tokens, o.webhookSecretFile)

	if err := secret.Add(tokens...); err != nil {
		logrus.WithError(err).Fatal("Error starting secrets agent.")
	}

	githubClient, err := o.github.GitHubClient(false)
	if err != nil {
		logger.WithError(err).Fatal("Error getting GitHub client.")
	}

	jiraClient, err := o.jira.Client()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to construct Jira Client")
	}

	serv := &server{
		config: func() *Config {
			o.mut.Lock()
			defer o.mut.Unlock()
			return o.config
		},
		ghc:             githubClient.WithFields(logger.Data).ForPlugin(PluginName),
		jc:              jiraClient.WithFields(logger.Data).ForPlugin(PluginName),
		prowConfigAgent: configAgent,
	}

	eventServer := githubeventserver.New(o.githubEventServerOptions, secret.GetTokenGenerator(o.webhookSecretFile), logger)
	eventServer.RegisterHandleIssueCommentEvent(serv.handleIssueComment)
	eventServer.RegisterHandlePullRequestEvent(serv.handlePullRequest)
	eventServer.RegisterHelpProvider(serv.helpProvider, logger)

	health := pjutil.NewHealth()
	health.ServeReady()

	interrupts.ListenAndServe(eventServer, time.Second*30)
	interrupts.WaitForGracefulShutdown()
}
