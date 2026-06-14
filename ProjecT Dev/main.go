package main

import (
	controllers "RedProject/controllers"
	"RedProject/database"
	routes "RedProject/routes"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

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
	flag.Parse()

	if *kill80 {
		killProcessOnPort("80")
	}
	if *kill8080 {
		killProcessOnPort("8080")
	}

	if !isPortAvailable("8080") {
		log.Fatalf("Le port 8080 est déjà utilisé. Utilisez --kill-port-8080 pour libérer le port, ou fermez le processus manuellement.")
	}

	database.InitDB()
	database.EnsureAdminAccount()

	controllers.Init()
	routes.InitRoutes()
	fs := http.FileServer(http.Dir("assets"))
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
