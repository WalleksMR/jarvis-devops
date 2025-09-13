// Configuration editor JavaScript
function configEditor(filename) {
  return {
    filename: filename,
    content: "",
    originalContent: "",
    hasChanges: false,
    loading: true,
    saving: false,
    validating: false,
    editor: null,
    validationResult: {
      show: false,
      valid: false,
      error: "",
    },
    notification: {
      show: false,
      type: "info",
      message: "",
    },

    async init() {
      await this.loadConfig();
      this.initializeEditor();
      this.setupKeyboardShortcuts();
    },

    async loadConfig() {
      this.loading = true;
      try {
        const response = await fetch(`/api/config/${this.filename}`);
        if (response.ok) {
          const data = await response.json();
          this.content = data.content;
          this.originalContent = data.content;
          this.hasChanges = false;
        } else {
          const error = await response.json();
          this.showNotification(
            "error",
            error.error || "Erro ao carregar arquivo"
          );
        }
      } catch (error) {
        this.showNotification("error", "Erro de conexão: " + error.message);
      } finally {
        this.loading = false;
      }
    },

    initializeEditor() {
      // Wait for next tick to ensure textarea is rendered
      this.$nextTick(() => {
        const textarea = document.getElementById("config-editor");
        
        // Check if editor is already initialized
        if (this.editor || !textarea || typeof CodeMirror === "undefined") {
          return;
        }

        this.editor = CodeMirror.fromTextArea(textarea, {
          mode: "nginx",
          theme: "material",
          lineNumbers: true,
          lineWrapping: true,
          indentUnit: 4,
          tabSize: 4,
          autoCloseBrackets: true,
          matchBrackets: true,
          showCursorWhenSelecting: true,
          styleActiveLine: true,
        });
        
        this.editor.setSize("100%", "600px");
        this.editor.setValue(this.content);

        this.editor.on("change", () => {
          this.content = this.editor.getValue();
          this.hasChanges = this.content !== this.originalContent;
        });
      });
    },

    setupKeyboardShortcuts() {
      document.addEventListener("keydown", (e) => {
        // Ctrl+S or Cmd+S to save
        if ((e.ctrlKey || e.metaKey) && e.key === "s") {
          e.preventDefault();
          if (this.hasChanges && !this.saving) {
            this.saveConfig();
          }
        }
      });
    },

    async saveConfig() {
      if (!this.hasChanges || this.saving) return;

      this.saving = true;
      try {
        const response = await fetch(`/api/config/${this.filename}`, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            content: this.content,
          }),
        });

        const data = await response.json();

        if (response.ok) {
          this.originalContent = this.content;
          this.hasChanges = false;
          this.showNotification(
            "success",
            data.message || "Arquivo salvo com sucesso"
          );
        } else {
          this.showNotification(
            "error",
            data.error || "Erro ao salvar arquivo"
          );
        }
      } catch (error) {
        this.showNotification("error", "Erro de conexão: " + error.message);
      } finally {
        this.saving = false;
      }
    },

    async validateConfig() {
      this.validating = true;
      this.validationResult.show = false;

      try {
        const response = await fetch("/api/validate", { method: "POST" });
        const data = await response.json();

        this.validationResult = {
          show: true,
          valid: data.valid,
          error: data.error || "",
        };

        if (data.valid) {
          this.showNotification("success", "Configuração válida!");
        } else {
          this.showNotification("error", "Configuração inválida");
        }
      } catch (error) {
        this.showNotification(
          "error",
          "Erro ao validar configuração: " + error.message
        );
      } finally {
        this.validating = false;
      }
    },

    async reloadConfig() {
      if (this.hasChanges) {
        if (
          !confirm(
            "Você tem alterações não salvas. Tem certeza que deseja recarregar o arquivo?"
          )
        ) {
          return;
        }
      }

      await this.loadConfig();
      
      // Update editor content if it exists
      if (this.editor) {
        this.editor.setValue(this.content);
        this.originalContent = this.content;
        this.hasChanges = false;
      }
      
      this.showNotification("info", "Arquivo recarregado");
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
  };
}
