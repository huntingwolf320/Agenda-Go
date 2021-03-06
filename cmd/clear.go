// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/Agenda-Go/entity"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete all meetings you sponsored.",
	Long:  "Usage: agenda clear.",
	Run: func(cmd *cobra.Command, args []string) {
		entity.ClearMeetings(entity.GetCurrentUser())
		log.Println("Clear completed.")
		entity.UpdateLib()
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
	entity.Init()

}
