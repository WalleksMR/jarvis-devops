import { Button } from "@/components/ui/button";
import { CardStatus } from "./card-status";
import { CardVersion } from "./card-version";
import { useState } from "react";
import { CardLogs } from "./card-logs";

export function Nginx() {
  const [nginxStatus, setNginxStatus] = useState({
    installed: true,
    running: true,
    configured: true,
    version: "1.0.0",
  });

  return (
    <>
      <div>
        <h3>Gerenciamento do Nginx</h3>
      </div>

      <div className="rounded-md bg-gray-100 p-10">
        <h3 className="mb-5 text-2xl font-bold">Status</h3>
        <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
          <CardStatus
            name="Instalado"
            status={nginxStatus.installed}
            statusMessage={nginxStatus.installed ? "Sim" : "Não"}
          />
          <CardStatus
            name="Executando"
            status={nginxStatus.running}
            statusMessage={nginxStatus.running ? "Sim" : "Não"}
          />
          <CardStatus
            name="Configuração"
            status={nginxStatus.configured}
            statusMessage={nginxStatus.configured ? "Válido" : "Inválido"}
          />
          <CardVersion name="Versão" version={nginxStatus.version} />
        </div>

        {!nginxStatus.installed && (
          <CardLogs
            title="Erro de instalação: "
            errorMessage="O Nginx não está instalado. Por favor, verifique o serviço."
            color="red"
          />
        )}
        {nginxStatus.installed && !nginxStatus.running && (
          <CardLogs errorMessage="O Nginx não está em execução. Por favor, verifique o serviço." />
        )}
        {nginxStatus.running && !nginxStatus.configured && (
          <CardLogs errorMessage="A configuração do Nginx está inválida. Por favor, verifique os arquivos de configuração." />
        )}
      </div>

      <div className="rounded-md bg-gray-100 p-10">
        <h3 className="mb-5 text-2xl font-bold">Ações</h3>
        <div className="flex space-x-2">
          <Button
            disabled={!nginxStatus.running}
            className={`${!nginxStatus.running ? "bg-gray-400" : "bg-emerald-500"} text-white`}
          >
            Validar Configuração
          </Button>
          <Button
            disabled={!nginxStatus.running}
            className={`text-white ${!nginxStatus.running ? "bg-gray-400" : "bg-blue-500"}`}
          >
            Recarregar
          </Button>

          <Button
            disabled={!nginxStatus.running}
            className={`text-white ${!nginxStatus.running ? "bg-gray-400" : "bg-yellow-500"}`}
          >
            Reiniciar
          </Button>
        </div>
      </div>
    </>
  );
}
