package main

import (
	"archive/tar"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func createOVA(files []string, ovaName string) {

	ovaFile, err := os.Create(ovaName + ".ova")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Create a new tar archive.
	tw := tar.NewWriter(ovaFile)

	// Add some files to the archive.
	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("%v", err)
		}
		defer file.Close()
		if stat, err := file.Stat(); err == nil {
			// now lets create the header as needed for this file within the tarball
			header := new(tar.Header)
			header.Name = filePath
			header.Size = stat.Size()
			header.Mode = int64(stat.Mode())
			header.ModTime = stat.ModTime()
			// write the header to the tarball archive
			if err := tw.WriteHeader(header); err != nil {
				log.Fatalf("%v", err)
			}
			// copy the file data to the tarball
			if _, err := io.Copy(tw, file); err != nil {
				log.Fatalf("%v", err)
			}
		}
	}
}

func writeOVF(ovfFileName string, ovfContent string) {
	ovaFile, err := os.Create(ovfFileName + ".ovf")
	if err != nil {
		log.Fatalf("%v", err)
	}
	_, err = ovaFile.WriteString(ovfContent)
	if err != nil {
		log.Fatalf("%v", err)
	}
	ovaFile.Sync()
}

func main() {
	invoked := filepath.Base(os.Args[0])

	network := flag.String("network", "", "The network label the VM will use")
	cpus := flag.String("cpus", "1", "Number of CPUs")
	mem := flag.String("mem", "1024", "Amount of memory in MB")

	flag.Usage = func() {
		fmt.Printf("USAGE: %s push vcenter [options] path \n\n", invoked)
		fmt.Printf("'path' specifies the full path of an image that will be pushed\n")
		fmt.Printf("Options:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	remArgs := flag.Args()
	if len(remArgs) == 0 {

		fmt.Printf("Please specify the path to the image to push\n")
		flag.Usage()
		os.Exit(1)
	}
	isoPath := remArgs[0]

	// Ensure an iso has been passed to the vCenter push Command
	if !strings.HasSuffix(isoPath, ".iso") {
		log.Fatalln("Please pass an \".iso\" file as the path")
	}

	envelope := newDMTFEnvelope()
	vmName := strings.TrimSuffix(path.Base(isoPath), ".iso")

	var net NetSection
	net.NetInfo = "List of Networks"
	net.Network.NetName = *network
	net.Network.NetDesc = "Default Network"
	envelope.Net = &net

	var newHardware VirtualHardware
	envelope.VM.VInfo = "A Virtual Machine"
	envelope.VM.VID = "vm"
	envelope.VM.VName = vmName
	envelope.VM.VOSSection.VOSID = "1"
	envelope.VM.VOSSection.VOSType = "*other3xLinux64Guest"
	envelope.VM.VOSSection.VOSInfo = "The kind of installed guest operating system"
	newHardware.VHWInfo = "Virtual hardware requirements"
	newHardware.VHWSystem.VHWInstanceID = "0"
	newHardware.VHWSystem.VHWSystemType = "vmx-11"
	newHardware.VHWSystem.VHWSystemID = vmName
	newHardware.VHWSystem.VHWSystemName = "Virtual Hardware Family"
	addMemoryToVM(&newHardware, *mem)
	addCPUtoVM(&newHardware, *cpus)
	addNicToVM(&newHardware, *network)
	ideController := addIDEControllerToVM(&newHardware)
	//	scsiController := addSCSIControllertoVM(&newHardware)

	fi, e := os.Stat(isoPath)
	if e != nil {
		os.Exit(1)
	}

	var newDiskSection DiskSection
	addCDToController(&newHardware, ideController, "file1")
	appendFilesToReferences(&envelope, path.Base(isoPath), "file1", fmt.Sprintf("%d", fi.Size()))

	// ISOs dont go in the disk section OBVIOUSLY :'''(
	//appdendDiskToDiskSection(&newDiskSection, "8", "byte * 2^30", "file1", "vmdisk1", "http://format", "0")

	newDiskSection.DiskInfo = "Virtual disk information"
	envelope.Disk = append(envelope.Disk, &newDiskSection)

	envelope.VM.VHardware = newHardware

	output, err := xml.MarshalIndent(envelope, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	var xmlOutput string
	xmlOutput = xml.Header + string(output)

	writeOVF(vmName, xmlOutput)
	createOVA([]string{vmName + ".ovf", isoPath}, vmName)
}
