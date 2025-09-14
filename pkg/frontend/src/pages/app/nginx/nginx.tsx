import { Button } from "@/components/ui/button";
import { CardStatus } from "./card-status";
import { CardVersion } from "./card-version";
import { useEffect, useState } from "react";
import { CardLogs } from "./card-logs";
import { useQuery } from "@tanstack/react-query";
import { getStatus } from "@/api/get-status";

export function Nginx() {
  const [nginx, setNginx] = useState<{
    installed: boolean | null;
    running: boolean | null;
    configured: boolean | null;
    version: string | null;
  }>({
    installed: null,
    running: null,
    configured: null,
    version: null,
  });

  const { data: nginxGetStatus, error } = useQuery({
    queryKey: ["nginx-get-status"],
    queryFn: getStatus,
  });

  useEffect(() => {
    if (nginxGetStatus) {
      setNginx({
        version: nginxGetStatus.version,
        installed: nginxGetStatus.is_installed,
        running: nginxGetStatus.is_running,
        configured: nginxGetStatus.config_valid,
      });
    }
  }, [nginxGetStatus]);

  return (
    <>
      <div>
        <h3 className="text-2xl font-medium">Dashboard</h3>
      </div>

      <div className="rounded-md bg-gray-100 p-10">
        <h3 className="mb-5 text-2xl font-light">Status</h3>
        <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
          <CardStatus
            name="Instalado"
            status={nginx.installed}
            statusMessage={nginx.installed ? "Sim" : "Não"}
          />
          <CardStatus
            name="Executando"
            status={nginx.running}
            statusMessage={nginx.running ? "Sim" : "Não"}
          />
          <CardStatus
            name="Configuração"
            status={nginx.configured}
            statusMessage={nginx.configured ? "Válido" : "Inválido"}
          />
          <CardVersion
            name="Versão"
            version={nginx.version || "Desconhecida"}
          />
        </div>
        {error && <CardLogs errorMessage={error.message} />}
        {!error ||
          (!nginx.installed && (
            <CardLogs
              title="Erro de instalação: "
              errorMessage={
                nginxGetStatus?.config_error
                  ? nginxGetStatus?.config_error
                  : "O Nginx não está instalado. Por favor, verifique o serviço."
              }
              color="red"
            />
          ))}

        {nginx.installed && !nginx.running && (
          <CardLogs
            errorMessage={
              nginxGetStatus?.config_error
                ? nginxGetStatus?.config_error
                : "O Nginx não está em execução. Por favor, verifique o serviço."
            }
          />
        )}
        {nginx.running && !nginx.configured && (
          <CardLogs
            errorMessage={
              nginxGetStatus?.config_error
                ? nginxGetStatus.config_error
                : "A configuração do Nginx está inválida. Por favor, verifique os arquivos de configuração."
            }
          />
        )}
      </div>

      <div className="rounded-md bg-gray-100 p-10">
        <h3 className="mb-5 text-2xl font-light">Ações</h3>
        <div className="flex space-x-2">
          <Button
            disabled={!nginx.running}
            className={`${!nginx.running ? "bg-gray-400" : "bg-emerald-500"} text-white`}
          >
            Validar Configuração
          </Button>
          <Button
            disabled={!nginx.running}
            className={`text-white ${!nginx.running ? "bg-gray-400" : "bg-blue-500"}`}
          >
            Recarregar
          </Button>

          <Button
            disabled={!nginx.running}
            className={`text-white ${!nginx.running ? "bg-gray-400" : "bg-yellow-500"}`}
          >
            Reiniciar
          </Button>
        </div>
      </div>
    </>
  );
}
