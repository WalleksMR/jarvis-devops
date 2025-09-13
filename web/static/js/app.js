// Main application JavaScript for nginx manager
function nginxManager() {
  return {
    status: {
      is_installed: false,
      is_running: false,
      version: "",
      config_valid: false,
      config_error: "",
    },
    configs: [],
    loading: {
      status: false,
      configs: false,
      validate: false,
      reload: false,
      restart: false,
      refresh: false,
    },
    notification: {
      show: false,
      type: "info", // success, error, warning, info
      message: "",
    },

    async init() {
      await this.refreshStatus();
      await this.loadConfigFiles();
    },

    async refreshStatus() {
      this.loading.refresh = true;
      try {
        const response = await fetch("/api/status");
        if (response.ok) {
          this.status = await response.json();
        } else {
          this.showNotification("error", "Erro ao carregar status do nginx");
        }
      } catch (error) {
        this.showNotification("error", "Erro de conexão: " + error.message);
      } finally {
        this.loading.refresh = false;
      }
    },

    async loadConfigFiles() {
      this.loading.configs = true;
      try {
        const response = await fetch("/api/configs");
        if (response.ok) {
          const data = await response.json();
          this.configs = data.configs || [];
        } else {
          this.showNotification(
            "error",
            "Erro ao carregar arquivos de configuração"
          );
        }
      } catch (error) {
        this.showNotification("error", "Erro de conexão: " + error.message);
      } finally {
        this.loading.configs = false;
      }
    },

    async validateConfig() {
      this.loading.validate = true;
      try {
        const response = await fetch("/api/validate", { method: "POST" });
        const data = await response.json();

        if (data.valid) {
          this.showNotification("success", "Configuração válida!");
          this.status.config_valid = true;
          this.status.config_error = "";
        } else {
          this.showNotification("error", "Configuração inválida");
          this.status.config_valid = false;
          this.status.config_error = data.error || "";
        }
      } catch (error) {
        this.showNotification(
          "error",
          "Erro ao validar configuração: " + error.message
        );
      } finally {
        this.loading.validate = false;
      }
    },

    async reloadNginx() {
      if (!this.status.config_valid) {
        this.showNotification(
          "warning",
          "Valide a configuração antes de recarregar"
        );
        return;
      }

      this.loading.reload = true;
      try {
        const response = await fetch("/api/reload", { method: "POST" });
        const data = await response.json();

        if (response.ok) {
          this.showNotification("success", data.message);
          await this.refreshStatus();
        } else {
          this.showNotification(
            "error",
            data.error || "Erro ao recarregar nginx"
          );
        }
      } catch (error) {
        this.showNotification("error", "Erro de conexão: " + error.message);
      } finally {
        this.loading.reload = false;
      }
    },

    async restartNginx() {
      if (!this.status.config_valid) {
        this.showNotification(
          "warning",
          "Valide a configuração antes de reiniciar"
        );
        return;
      }

      if (
        !confirm(
          "Tem certeza que deseja reiniciar o nginx? Isso pode interromper temporariamente o serviço."
        )
      ) {
        return;
      }

      this.loading.restart = true;
      try {
        const response = await fetch("/api/restart", { method: "POST" });
        const data = await response.json();

        if (response.ok) {
          this.showNotification("success", data.message);
          await this.refreshStatus();
        } else {
          this.showNotification(
            "error",
            data.error || "Erro ao reiniciar nginx"
          );
        }
      } catch (error) {
        this.showNotification("error", "Erro de conexão: " + error.message);
      } finally {
        this.loading.restart = false;
      }
    },

    editConfig(filename) {
      window.location.href = `/editor/${filename}`;
    },

    showNotification(type, message) {
      this.notification = {
        show: true,
        type: type,
        message: message,
      };

      // Auto-hide after 5 seconds
      setTimeout(() => {
        this.notification.show = false;
      }, 5000);
    },

    formatFileSize(bytes) {
      if (bytes === 0) return "0 Bytes";
      const k = 1024;
      const sizes = ["Bytes", "KB", "MB", "GB"];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
    },

    formatDate(dateString) {
      const date = new Date(dateString);
      return (
        date.toLocaleDateString("pt-BR") +
        " " +
        date.toLocaleTimeString("pt-BR")
      );
    },
  };
}
