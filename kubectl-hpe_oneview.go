package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/HewlettPackard/oneview-golang/utils"
	"github.com/mpolski/oneview-golang-temp/ov"
)

var (
	addnode      *int
	capacity     *bool
	h            *bool
	help         *bool
	compute      *bool
	removenode   *int
	storage 	 *bool
	template     *bool
	templates    *bool

	serverHardwareTypeURI         string
	serverProfileTemplateNameURI  string
	mServerHardware               map[string]utils.Nstring
	mServerHardwareType           map[string]string
	mServerHardwareListNoProfiles map[string]string
	mTemplates                    map[string]string

	clientOV *ov.OVClient

	spt = os.Getenv("OV_PROFILETEMPLATE")

	ovc = clientOV.NewOVClient(
		os.Getenv("OV_USERNAME"),
		os.Getenv("OV_PASSWORD"),
		os.Getenv("OV_AUTHLOGINDOMAIN"),
		os.Getenv("OV_ENDPOINT"),
		false,
		1000,
		"*")
)

func init() {

	//define flags
	addnode = flag.Int("addnode", 0, "Future use!")
	capacity = flag.Bool("capacity", false, "Lists of servers that can be used for Kubernetes cluster expansion.")
	help = flag.Bool("help", false, "Usage: [OPTION].. e.g. -compute")
	h = flag.Bool("h", false, "Usage: [OPTION].. e.g. -compute")
	compute = flag.Bool("compute", false, "Lists information about servers currently used in a cluster.") // - should be subcommand of spt
	removenode = flag.Int("removenode", 0, "Future use!")
	storage = flag.Bool("storage", false, "List disks attached to compute modules, either local or from disk shelf DS3940")
	templates = flag.Bool("templates", false, "Lists available templates")

	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Printf("\nUsage %s [options]\n", os.Args[0])
		fmt.Println("\nOptions:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {

	mServerHardware = make(map[string]utils.Nstring)
	mServerHardwareType = make(map[string]string)
	mServerHardwareListNoProfiles = make(map[string]string)
	mTemplates = make(map[string]string)

	//writer, minwidth, tabwidth, padding int, padchar byte, flags uint)
	w := tabwriter.NewWriter(os.Stdout, 20, 16, 2, ' ', 0)
	t := tabwriter.NewWriter(os.Stdout, 20, 16, 2, ' ', 0)

	//COMPUTE
	if *compute == true {
		comp, err := ovc.GetProfileTemplateByName(spt)
		if err != nil {
			fmt.Println("Server Profile Template Retrieval By Name Failed: ", err)
		} else {
			//define query filter
			f := "serverProfileTemplateUri matches '" + comp.URI.String() + "'"
			//Get profiles for the given template
			prof, err := ovc.GetProfiles("", "", f, "", "")
			if err != nil {
				fmt.Println("Server Profiles Retrieval by URI Failed: ", err)
			} else {
				//put profile names created out of that template and their corresponding serverHardware URIs into a map
				for i := 0; i < len(prof.Members); i++ {
					mServerHardware[prof.Members[i].Name] = prof.Members[i].ServerHardwareURI
				}
			}
			if len(mServerHardware) > 0 {
				fmt.Fprintln(w, "NAME\tvCPU\tRAM[GB]\tSTATUS\tPOWER STATE\tMODEL\tLOCATION [ENCLOSURE, BAY]\t")
				for _, v := range mServerHardware {
					consumedComp, err := ovc.GetServerHardwareByUri(v)
					if err != nil {
						fmt.Println("Server Hardware List for Status Retrival Failed: ", err)
					} else {
						hostName := consumedComp.ServerName
						vCPU := consumedComp.ProcessorCount * consumedComp.ProcessorCoreCount * 2
						memory := consumedComp.MemoryMb / 1024
						model := consumedComp.Model
						status := consumedComp.Status
						power := consumedComp.PowerState
						locationName := consumedComp.Name
						fmt.Fprintln(w, hostName, "\t", vCPU, "\t", memory, "\t", status, "\t", power, "\t", model, "\t", locationName, "\t") //needs servers.serverName - add to SDK
						w.Flush()
					}
				}
				fmt.Println("\t\nTOTAL NODES:  \t", len(mServerHardware), "\n")
				fmt.Println("Servers deployed using template:", spt, "\n")
			} else {
				fmt.Println("No servers running using template: ", spt, "\n")
			}
		}
	}

	//storage
	if *storage == true {
		stor, err := ovc.GetProfileTemplateByName(spt)
		if err != nil {
			fmt.Println("Server Profile Template Retrieval By Name Failed: ", err)
		} else {
			//define query filter
			f := "serverProfileTemplateUri matches '" + stor.URI.String() + "'"
			//Get profiles for the given template
			prof, err := ovc.GetProfiles("", "", f, "", "")
			if err != nil {
				fmt.Println("Server Profiles Retrieval by URI Failed: ", err)
			} else {
				//put profile names created out of that template and their corresponding serverHardware URIs into a map
				for i := 0; i < len(prof.Members); i++ {
					mServerHardware[prof.Members[i].Name] = prof.Members[i].ServerHardwareURI

				}
			}
			fmt.Println("\t\nLOCAL STORAGE - DISKS IN COMPUTE NODES OR DY3940 (P416):\n")
			if len(mServerHardware) > 0 {
				fmt.Fprintln(w, "  NAME\tCONTROLLER / STATUS\t\tDISK NO.\tCAPACITY[GB]\tINTERFACE\tMEDIA\tDISK HEALTH\tSERIAL NO.\tMODEL")

				var hostname string

				for _, v := range mServerHardware {
					consumedStor, err := ovc.GetServerHardwareByUri(v)
					if err != nil {
						fmt.Println("Server Hardware List for Status Retrival Failed: ", err)
					} else {
						hostname = consumedStor.ServerName
						storage, err := ovc.GetServerLocalStorageByUri(v)
						if err != nil {
							fmt.Println("Storage List for Server Hardware Retrival Failed: ", err)
						} else {
							data := storage.Data
							fmt.Fprintln(w, " ", hostname)
							for i := 0; i < len(data); i++ {
								ctrl := data[i].Model
								ctrlHealth := data[i].Status.Health
								fmt.Fprintln(w, "\t", ctrl, "/", ctrlHealth)
								for p := 0; p < len(data[i].PhysicalDrives); p++ {
									diskNo := p
									capacityGB := int(float64(data[i].PhysicalDrives[p].CapacityMiB) * 0.001048576) //MiB to GB conversion
									interfaceType := data[i].PhysicalDrives[p].InterfaceType
									mediaType := data[i].PhysicalDrives[p].MediaType
									model := data[i].PhysicalDrives[p].Model
									serialNumber := data[i].PhysicalDrives[p].SerialNumber
									diskHealth := data[i].PhysicalDrives[p].Status.Health
									fmt.Fprintln(w, "\t", "\t", "\t", diskNo, "\t", capacityGB, "\t", interfaceType, "\t", mediaType, "\t", diskHealth, "\t", serialNumber, "\t", model) //needs servers.serverName - add to SDK
									w.Flush()
								}
							}
						}
					}
				}
			} else {
				fmt.Println("No servers running using template: ", spt, "\n")
			}
		}
	}

	//CAPACITY
	if *capacity == true {
		cap, err := ovc.GetProfileTemplateByName(spt)
		if err != nil {
			fmt.Println("Server Profile Template Retrieval By Name Failed: ", err)
		} else {
			//first, search for those that match serverHardwareType
			f := []string{"serverHardwareTypeURI matches '" + cap.ServerHardwareTypeURI.String() + "'", "state matches Unmanaged"}
			c, err := ovc.GetServerHardwareList(f, "")
			if err != nil {
				fmt.Println("Server Hardware List Retrival Failed: ", err)
			} else {
				//put all servers of the same serverHardwareType into a map for further filtering
				for i := 0; i < len(c.Members); i++ {
					mServerHardwareType[c.Members[i].Name] = c.Members[i].URI.String()
				}
			}
			fmt.Println("\nCAPACITY:\t", len(mServerHardwareType), "\n")
			if len(mServerHardwareType) > 0 {
				fmt.Println("Available servers to deploy profile:", spt, "\n")
				fmt.Fprintln(w, "  Location\tvCPU\tRAM[GB]\tStatus\tPower State\t")
				fmt.Fprintln(w, "  ----------\t----------\t----------\t----------\t----------\t")
				for k := range mServerHardwareType {
					servers, err := ovc.GetServerHardwareByName(k)
					if err != nil {
						fmt.Println("Server Hardware List for Capacity Retrival Failed: ", err)
					} else {
						vCPU := servers.ProcessorCount * servers.ProcessorCoreCount * 2
						memory := servers.MemoryMb / 1024
						fmt.Fprintln(w, "  ", servers.Name, "\t", vCPU, "\t", memory, "\t", servers.Status, "\t", servers.PowerState, "\t") //needs servers.serverName - add to SDK
						w.Flush()
					}
				}
			} else {
				fmt.Println("No more servers available for template: ", spt, "\n Add more servers.\n")
			}
		}
	}

	//HELP or H
	if *help == true || *h == true {
		fmt.Println("\nUsage \"kubectl hpe oneview\" [options]\n")
		fmt.Println("\n[Options]\n")
		fmt.Println(" -addnode [int]		-Future use!")
		fmt.Println(" -capacity		-Lists of servers that can be used for Kubernetes cluster extension.")
		fmt.Println(" -compute		-Lists information about servers currently used in a cluster.")
		fmt.Println(" -removenode [int]	-Future use!")
		fmt.Println(" -storage		-Lists nodes and their attached storage details (disks and controllers) either onboard or from DY3940 disk shelf.")
		fmt.Println(" -templates		-Lists ready templates. Not all have deployment plans, confirm before use.")
		fmt.Println("\n")
	}

	//TEMPLATES
	if *templates == true {
		temp, err := ovc.GetProfileTemplates("", "", "", "", "")
		if err != nil {
			fmt.Println("Server Profile Templates Retrieval Failed: ", err)
		} else {
			for i := 0; i < len(temp.Members); i++ {
				mTemplates[temp.Members[i].Name] = temp.Members[i].ServerHardwareTypeURI.String()
			}
		}
		fmt.Println("\tTEMPLATES:  \t", len(mTemplates), "\n")

		if len(mTemplates) > 0 {
			fmt.Println("Availble templates to choose from:\n")
			fmt.Fprintln(t, "  Template Name\tNo. of compatible Srvs\t")
			fmt.Fprintln(t, "  --------------\t--------------\t")

			for k, v := range mTemplates {
				f := []string{"serverHardwareTypeURI matches '" + v + "'"}
				templates, err := ovc.GetServerHardwareList(f, "")
				if err != nil {
					fmt.Println("Server Hardware Type List by Template Failed: ", err)
				} else {
					fmt.Fprintln(t, "  ", k, "\t", templates.Count, "\t")
					t.Flush()
				}
			}
			fmt.Println("\n")
		} else {
			fmt.Println("No Server Profiles Templates available\n")
		}
	}

	//ADDNODE
	if *addnode > 0 {
		fmt.Println("\n NOT YET IMPLEMENTED \n")
	}

	//REMOVNODE
	if *removenode > 0 {
		fmt.Println("\n NOT YET IMPLEMENTED \n")
	}
}
