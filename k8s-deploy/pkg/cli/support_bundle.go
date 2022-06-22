/*
 * Copyright 2019-2021 VMware, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package cli

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	client "github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/k8sclient"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/utils/supportbundle"
	"github.com/urfave/cli/v2"
)

var (
	bundler *supportbundle.Bundler
)

func SupportBundleCommand() *cli.Command {
	return &cli.Command{
		Name:  "supportbundle",
		Flags: []cli.Flag{},
		Subcommands: []*cli.Command{
			SupportBundleCollectCommand(),
			SupportBundlePackCommand(),
		},
		Usage: "Collect data to perform troubleshooting",
	}
}

func SupportBundleCollectCommand() *cli.Command {
	return &cli.Command{
		Name: "collect",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "tail",
				Aliases: []string{"t"},
				Value:   200,
				Usage:   "Specify how many rows to record.",
			},
			&cli.StringFlag{
				Name:    "kubeconfig",
				Aliases: []string{"c"},
				Value:   "",
				Usage:   "Specify the kubeconfig.",
			},
			&cli.StringFlag{
				Name:  "collectDir",
				Value: "./kubefate-supportbundle",
				Usage: "Specify where to save temporary data.",
			},
			&cli.StringFlag{
				Name:  "packDir",
				Value: "./",
				Usage: "Specify where to pack zip file.",
			},
		},
		Usage: "Collect data",
		Action: func(c *cli.Context) error {
			tail := c.Int("tail")
			kubeconfig := c.String("kubeconfig")
			if kubeconfig == "" {
				kubeconfig = client.GetKubeconfig()
			}
			collectDir := c.String("collectDir")
			packDir := c.String("packDir")
			if err := initBundler(kubeconfig, collectDir, packDir); err != nil {
				return err
			}
			namespaceList, err := bundler.Client.GetNamespaceList()
			if err != nil {
				return err
			}
			var namespaces []string
			namespacePrompt := &survey.MultiSelect{
				Message: "Choose namespaces to troubleShooting:",
				Options: client.NamespaceListToNames(namespaceList),
				Default: []string{"kube-fate"},
			}
			survey.AskOne(namespacePrompt, &namespaces)

			for _, namespace := range namespaces {
				err := CollectInNamespace(bundler, namespace, tail)
				if err != nil {
					fmt.Println(err)
				}
			}
			fmt.Printf("You can attach additional files about this issue into %s\n", collectDir)
			return nil
		},
	}
}

func SupportBundlePackCommand() *cli.Command {
	return &cli.Command{
		Name: "pack",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "tail",
				Aliases: []string{"t"},
				Value:   200,
				Usage:   "Specify how many rows to record.",
			},
			&cli.StringFlag{
				Name:    "kubeconfig",
				Aliases: []string{"c"},
				Value:   "",
				Usage:   "Specify the kubeconfig.",
			},
			&cli.StringFlag{
				Name:  "collectDir",
				Value: "./kubefate-supportbundle",
				Usage: "Specify where to save temporary data.",
			},
			&cli.StringFlag{
				Name:  "packDir",
				Value: "./",
				Usage: "Specify where to pack zip file.",
			},
		},
		Usage: "Pack data",
		Action: func(c *cli.Context) error {
			kubeconfig := c.String("kubeconfig")
			if kubeconfig == "" {
				kubeconfig = client.GetKubeconfig()
			}
			collectDir := c.String("collectDir")
			packDir := c.String("packDir")
			if err := initBundler(kubeconfig, collectDir, packDir); err != nil {
				return err
			}

			privacyAck := true
			privacyPrompt := &survey.Confirm{
				Message: "This program collects some data, please ensure that there is no any privacy issue.\n Enter y if you're willing to share these collected data with us.",
			}
			survey.AskOne(privacyPrompt, &privacyAck)
			if !privacyAck {
				return errors.New("please acknowledge the privacy agreement")
			}
			return bundler.Pack()
		},
	}
}

func CollectInNamespace(b *supportbundle.Bundler, namespace string, tail int) (err error) {
	podList, err := b.Client.GetPodList(namespace)
	if err != nil {
		return err
	}
	var pods []string
	podsPrompt := &survey.MultiSelect{
		Message: fmt.Sprintf("In namespace %s, choose pods to troubleShooting:", namespace),
		Options: client.PodListToNames(podList),
	}
	survey.AskOne(podsPrompt, &pods)
	return b.CollectNamespace(namespace, tail, pods...)
}

func initBundler(kubeconfig, collectDir, packDir string) (err error) {
	bundler, err = supportbundle.NewBundler(kubeconfig, collectDir, packDir)
	if err != nil {
		return
	}
	if bundler == nil {
		err = errors.New("supportbundler not initialized sucessfully")
	}
	return
}
