package cli

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"os"
	"time"
)

func initClis() {
	Register( "echo", &CliFuncPack{ echo, "Echo text", "basic" } )
	Register( "help", &CliFuncPack{ help, "List all commands and descriptions", "basic" } )
	Register( "stop", &CliFuncPack{ stop, "Abort application", "basic" } )
	Register( "hds", &CliFuncPack{ hds, "Print system hardware information", "sys" } )
}

func echo( ts []string ) error {
	str := ""
	for _, t := range ts {
		str += t + " "
	}
	println( str )
	return nil
}

func help( ts []string ) error {
	println("===")
	for g, _ := range groups {
		println( "[" + g + "]" )
		for n, p := range CliFuncs {
			if p.Group == g {
				println("\t"+ n + " - " + p.Desc)
			}
		}
	}
	println("===")
	return nil
}

func stop( ts []string ) error {
	print("server stopped\t")
	os.Exit(0)
	return nil
}

func hds( ts []string ) error {
	println("hardware monitor:")
	println("Total RAM: ", getMemPercent() )
	return nil
}

func getCpuPercent() float64 {
	percent, _:= cpu.Percent(time.Second, false)
	return percent[0]
}

func getMemPercent()uint64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.Total
}