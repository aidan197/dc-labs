package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type packages struct {
	name        string
	installDate string
	updateDate  string
	nUpdates    int
	removalDate string
}
type packageManager struct {
	installed int
	remove    int
	upgrade   int
	current   int
}

func main() {

	pacman := packageManager{installed: 0, remove: 0, upgrade: 0, current: 0}
	var txtlines []string
	var packageList []packages

	file, err := os.Open("pacman.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	for _, eachline := range txtlines {
		values := strings.Fields(eachline)
		if values[2] == "[ALPM]" {
			if values[3] == "installed" {
				packageList = append(packageList, pacman.addP(values))
			} else if values[3] == "upgraded" {
				index := 0
				for i := 0; i <= len(packageList); i++ {
					if packageList[i].name == values[4] {
						index = i
						break
					}
				}
				packageList[index].nUpdates++
				packageList[index].updateDate = pacman.upgradeP(values)

			} else if values[3] == "removed" {
				index := 0
				for i := 0; i <= len(packageList); i++ {
					if packageList[i].name == values[4] {
						index = i
						break
					}
				}
				packageList[index].removalDate = pacman.removeP(values)
			} else if values[3] == "reinstalled" {
				index := 0
				for i := 0; i <= len(packageList); i++ {
					if packageList[i].name == values[4] {
						index = i
						break
					}
				}
				pacman.readdP()
				packageList[index].readdP(values)
			}
		}
	}

	f, err := os.Create("packages_report.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString(pacman.toString())

	for _, p := range packageList {
		f.WriteString(p.toString())
	}

}

func (pacman *packageManager) addP(txtLine []string) packages {
	date := []string{txtLine[0], txtLine[1]}
	packageData := packages{installDate: strings.Join(date, " "), name: txtLine[4], nUpdates: 0, removalDate: "-", updateDate: "-"}
	pacman.installed++
	pacman.current++
	return packageData
}
func (pack *packages) readdP(txtLine []string) {
	date := []string{txtLine[0], txtLine[1]}
	pack.installDate = strings.Join(date, " ")
	pack.nUpdates = 0
	pack.removalDate = "-"
	pack.updateDate = "-"
}
func (pacman *packageManager) readdP() {
	pacman.current++
}
func (pacman *packageManager) upgradeP(txtLine []string) string {
	date := []string{txtLine[0], txtLine[1]}
	pacman.upgrade++
	return strings.Join(date, " ")
}
func (pacman *packageManager) removeP(txtLine []string) string {
	date := []string{txtLine[0], txtLine[1]}
	pacman.remove++
	pacman.current--
	return strings.Join(date, " ")
}

func (pacman packageManager) toString() string {
	res := "Pacman Packages Report" + "\n"
	res += "----------------------" + "\n"
	res += strings.Join([]string{"- Installed packages :", strconv.Itoa(pacman.installed)}, " ") + "\n"
	res += strings.Join([]string{"- Removed packages :", strconv.Itoa(pacman.remove)}, " ") + "\n"
	res += strings.Join([]string{"- Upgraded packages :", strconv.Itoa(pacman.upgrade)}, " ") + "\n"
	res += strings.Join([]string{"- Current packages :", strconv.Itoa(pacman.current)}, " ") + "\n"
	res += "" + "\n"
	res += "List of packages" + "\n"
	res += "----------------" + "\n"
	return res
}

func (pack packages) toString() string {
	res := ""
	res += "- Package Name        : " + pack.name + "\n"
	res += "  - Last update date  : " + pack.updateDate + "\n"
	res += strings.Join([]string{"  - How many updates  : ", strconv.Itoa(pack.nUpdates)}, " ") + "\n"
	res += "  - Removal date      : " + pack.removalDate + "\n"
	return res
}
