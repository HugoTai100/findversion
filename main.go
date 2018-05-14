package main

import (
	"context"
	"fmt"
	"os"
	"bufio"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	if(len(releases) == 0){
		return versionSlice
	}
	semver.Sort(releases)
	for i := len(releases)-1; i >= 0; i-- {
	//should not use prereleases for productive use
		if(releases[i].PreRelease == "" && releases[i].Compare(*minVersion) >= 0) {
			releases[i].Metadata = "" //Build metadata SHOULD be ignored when determining version precedence.
			if(len(versionSlice) == 0){
				versionSlice = append(versionSlice, releases[i])
			}else if(versionSlice[len(versionSlice)-1].Major != releases[i].Major || versionSlice[len(versionSlice)-1].Minor != releases[i].Minor){
				 versionSlice = append(versionSlice,releases[i])
			}
		}
	}
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func getVersionArray(username, repo string)([]*semver.Version) {
	// Github
	errorVer := semver.New("666.666.666") // return this version for error
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	releases, _, err := client.Repositories.ListReleases(ctx, username, repo, opt)
	if err != nil {
		//fmt.Println(err)
		return []*semver.Version{errorVer}
	}
	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		defer func() {
			if r := recover(); r != nil {
				//fmt.Println(r) //Print error
				allReleases = []*semver.Version{errorVer}
			}
		}()
		allReleases[i] = semver.New(versionString)
	}
	return allReleases
}
func mapVersions(releases []*semver.Version, versionSlice []*semver.Version,  minVersion *semver.Version)[]*semver.Version {
	return LatestVersions(releases, minVersion)
}
func main() {
	if(len(os.Args) > 1){
		inFile, _ := os.Open(os.Args[1])
		defer inFile.Close()
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)
		scanner.Scan()
		for scanner.Scan() {
			var line1 string = scanner.Text()
			i1 := strings.Index(line1, "/")
			i2 := strings.Index(line1, ",")
			minVersion := semver.New(line1[i2+1:])
			allReleases := getVersionArray(line1[:i1], line1[i1+1:i2])
			versionSlice := LatestVersions(allReleases, minVersion)
			fmt.Printf("latest versions of %s/%s: %s\n",line1[:i1],line1[i1+1:i2], versionSlice)
		}
	}
}
