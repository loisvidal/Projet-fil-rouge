package main

import (
	controllers "RedProject/controllers"
	"RedProject/database"
	routes "RedProject/routes"
	"context"
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func init() {
	mime.AddExtensionType(".webp", "image/webp")
	mime.AddExtensionType(".jpeg", "image/jpeg")
	mime.AddExtensionType(".jpg", "image/jpeg")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".png", "image/png")
	mime.AddExtensionType(".svg", "image/svg+xml")
}

func findAssetsDir() string {
	cwd, _ := os.Getwd()
	cwdAssets := filepath.Join(cwd, "assets")
	if info, err := os.Stat(cwdAssets); err == nil && info.IsDir() {
		return cwdAssets
	}
	exePath, err := os.Executable()
	if err == nil {
		exeAssets := filepath.Join(filepath.Dir(exePath), "assets")
		if info, err := os.Stat(exeAssets); err == nil && info.IsDir() {
			return exeAssets
		}
	}
	return cwdAssets
}

func isWindowsAdmin() bool {
	if runtime.GOOS != "windows" {
		return true
	}
	cmd := exec.Command("net", "session")
	return cmd.Run() == nil
}

func elevateIfNeeded(portFlag string) {
	if isWindowsAdmin() {
		return
	}
	exe, err := os.Executable()
	if err != nil {
		log.Printf("Impossible de déterminer l'exécutable: %v", err)
		return
	}
	log.Printf("Tentative d'élévation des privilèges pour %s...\n", portFlag)
	psCmd := fmt.Sprintf("Start-Process -FilePath '%s' -ArgumentList '%s' -Verb RunAs -WindowStyle Hidden", exe, portFlag)
	cmd := exec.Command("powershell", "-Command", psCmd)
	if err := cmd.Start(); err != nil {
		log.Printf("Erreur élévation: %v", err)
		return
	}
	log.Println("Processus administrateur lancé. Fermeture de celui-ci.")
	os.Exit(0)
}

func killProcessOnPort(port string) {
	cmd := exec.Command("netstat", "-ano")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Erreur netstat: %v", err)
		return
	}

	lines := strings.Split(string(out), "\n")
	pids := make(map[string]bool)

	for _, line := range lines {
		if !strings.Contains(line, "LISTENING") {
			continue
		}
		if !strings.Contains(line, ":"+port+" ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			pid := fields[len(fields)-1]
			if pid != "0" && pid != "" {
				pids[pid] = true
			}
		}
	}

	if len(pids) == 0 {
		log.Printf("Aucun processus trouvé sur le port %s.", port)
		return
	}

	for pid := range pids {
		log.Printf("Arrêt du processus %s (port %s)...", pid, port)
		killCmd := exec.Command("taskkill", "/PID", pid, "/F")
		if output, err := killCmd.CombinedOutput(); err != nil {
			log.Printf("Erreur arrêt processus %s: %v - %s", pid, err, string(output))
		} else {
			log.Printf("Processus %s arrêté avec succès.", pid)
		}
	}
}

func isPortAvailable(port string) bool {
	cmd := exec.Command("netstat", "-ano")
	out, err := cmd.Output()
	if err != nil {
		return true
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "LISTENING") && strings.Contains(line, ":"+port+" ") {
			return false
		}
	}
	return true
}

func main() {
	kill80 := flag.Bool("kill-port-80", false, "Tue le processus utilisant le port 80")
	kill8080 := flag.Bool("kill-port-8080", false, "Tue le processus utilisant le port 8080")
	seed := flag.Bool("seed", false, "Réinitialise la base et insère des données de démonstration")
	flag.Parse()

	if *seed {
		database.ResetAndSeed()
		return
	}

	if *kill80 {
		elevateIfNeeded("--kill-port-80")
		killProcessOnPort("80")
	}
	if *kill8080 {
		elevateIfNeeded("--kill-port-8080")
		killProcessOnPort("8080")
	}

	if !isPortAvailable("8080") {
		log.Fatalf("Le port 8080 est déjà utilisé. Utilisez --kill-port-8080 pour libérer le port, ou fermez le processus manuellement.")
	}

	database.InitDB()
	database.EnsureAdminAccount()
	database.EnsurePropertyImages()

	controllers.Init()
	routes.InitRoutes()
	assetsDir := findAssetsDir()
	log.Printf("Assets servis depuis: %s", assetsDir)
	fs := http.FileServer(http.Dir(assetsDir))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Println("Serveur lancé sur http://localhost:8080/red_project/home")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur serveur : %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Arrêt du serveur en cours...")

	cleanupExe()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erreur lors de l'arrêt : %v", err)
	}

	log.Println("Serveur arrêté avec succès.")
}

func cleanupExe() {
	selfName := "RedProject.exe"
	if _, err := os.Stat(selfName); err == nil {
		rename := fmt.Sprintf("RedProject_%d.exe", time.Now().Unix())
		if err := os.Rename(selfName, rename); err == nil {
			log.Printf("Exécutable précédent renommé en %s", rename)
		}
	}
}
