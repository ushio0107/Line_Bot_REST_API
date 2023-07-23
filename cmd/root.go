/*
Copyright © 2023 Leung Yan Tung <leungyantung0107@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"m800_homework/api"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	lineCfgFile    = "./config/line_config.yml"
	mongodbCfgFile = "./config/mongodb_config.yml"
	serverCfgFile  = "./config/server_config.yml"
)

var (
	lineCfg   string
	dbCfg     string
	serverCfg string
)

var cfg = api.Config{
	LineCfg:   lineCfgFile,
	DBCfg:     mongodbCfgFile,
	ServerCfg: serverCfgFile,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "m800_homework",
	Short: "A simple RESTful API",
	Long: `This is a simple RESTful API which designed for
	Line Bot`,
	Run: func(cmd *cobra.Command, args []string) {
		exec.Command("make", "docker_install").Run()
		fmt.Print("Run root")
		a, err := api.NewServer(&cfg)
		if err != nil {
			log.Fatal("Fail: ", err)
		}
		a.Run()
		defer api.CloseDatabase(a.MongoClient)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&cfg.LineCfg, "lineCfg", lineCfgFile, "")
	rootCmd.Flags().StringVar(&cfg.ServerCfg, "serverCfg", serverCfgFile, "")
	rootCmd.Flags().StringVar(&dbCfg, "dbCfg", mongodbCfgFile, "")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
