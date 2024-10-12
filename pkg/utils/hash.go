package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ashwiniag/goKakashi/notifier"
	"io/ioutil"
	"os"
	"strings"
)

type HashEntry struct {
	Image           string              `json:"image"`           // Image name
	Tag             string              `json:"tag"`             // Image tag
	Vulnerabilities []VulnerabilityData `json:"vulnerabilities"` // Detailed vulnerability data
	Hash            string              `json:"hash"`            // Generated hash for the entry
}

type VulnerabilityData struct {
	VulnerabilityID  string `json:"vulnerability_id"`  // The CVE or vulnerability ID
	Severity         string `json:"severity"`          // Severity level (e.g., Critical, High)
	InstalledVersion string `json:"installed_version"` // Version of the package installed
	FixedVersion     string `json:"fixed_version"`     // Version where the vulnerability is fixed (if available)
}

// Path to store hash JSON
const hashFilePath = "./hashes.json"

// ConvertVulnerabilities converts []notifier.Vulnerability to []VulnerabilityData for hash storage
func ConvertVulnerabilities(vulnerabilities []notifier.Vulnerability) ([]VulnerabilityData, []string) {
	var vulnerabilityData []VulnerabilityData
	var vulnerabilityEntries []string
	for _, v := range vulnerabilities {
		data := VulnerabilityData{
			VulnerabilityID:  v.VulnerabilityID,
			Severity:         v.Severity,
			InstalledVersion: v.InstalledVersion,
			FixedVersion:     v.FixedVersion,
		}
		vulnerabilityData = append(vulnerabilityData, data)
		entry := fmt.Sprintf("%s_%s_%s_%s", data.VulnerabilityID, data.Severity, data.InstalledVersion, data.FixedVersion)
		vulnerabilityEntries = append(vulnerabilityEntries, entry)
	}
	return vulnerabilityData, vulnerabilityEntries
}

// GenerateHash creates a unique hash from image name, tag, and vulnerabilities
func GenerateHash(image string, tag string, vulnerabilities []string) string {
	data := fmt.Sprintf("%s:%s_%s", image, tag, strings.Join(vulnerabilities, "_")) // image/tag+detected vulenrabilties
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// SaveHashToFile saves the hash to the JSON file
func SaveHashToFile(filePath string, entry HashEntry) error {
	var entries []HashEntry

	file, err := ioutil.ReadFile(filePath)
	if err == nil {
		json.Unmarshal(file, &entries)
	}

	entries = append(entries, entry)

	fileContent, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, fileContent, 0644)
}

// HashExists checks if the hash already exists in the JSON file
func HashExists(filePath string, hash string) (bool, error) {
	var entries []HashEntry

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // No file yet, so no hash exists
		}
		return false, err
	}

	json.Unmarshal(file, &entries)

	for _, entry := range entries {
		if entry.Hash == hash {
			return true, nil
		}
	}

	return false, nil
}
