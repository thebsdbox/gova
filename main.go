package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"
)

func main() {

	var envelope Envelope

	envelope.BuildID = "build123"

	var net NetSection
	net.NetInfo = "List of Networks"
	net.Network.NetName = "VM Network"
	net.Network.NetDesc = "Default Network"
	envelope.Net = &net

	var newHardware VirtualHardware
	newHardware.VHWSystem.VHWInstanceID = "0"
	newHardware.VHWSystem.VHWSystemType = "vmx-11"
	newHardware.VHWSystem.VHWSystemID = "Other Linux 3.x kernel 64-bit"
	newHardware.VHWSystem.VHWSystemName = "Virtual Hardware Family"
	addMemoryToVM(&newHardware, "2048")
	addCPUtoVM(&newHardware, "4")

	ideController := addIDEControllerToVM(&newHardware)
	//	scsiController := addSCSIControllertoVM(&newHardware)

	cdFilePath := "./vcenter.iso"

	fi, e := os.Stat(cdFilePath)
	if e != nil {
		os.Exit(1)
	}
	var newDiskSection DiskSection
	addCDToController(&newHardware, ideController, cdFilePath)
	appendFilesToReferences(&envelope, path.Base(cdFilePath), "file1", fmt.Sprintf("%d", fi.Size()))

	// ISOs dont go in the disk section OBVIOUSLY :'''(
	//appdendDiskToDiskSection(&newDiskSection, "8", "byte * 2^30", "file1", "vmdisk1", "http://format", "0")

	newDiskSection.DiskInfo = "Virtual disk information"
	envelope.Disk = append(envelope.Disk, &newDiskSection)

	envelope.VM.VHardware.Hardware = newHardware

	output, err := xml.MarshalIndent(envelope, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	var xmlOutput string
	xmlOutput = xml.Header + string(output)
	fmt.Println(xmlOutput)
}
