type CardLogsProps = {
  title?: string;
  errorMessage?: string;
  color?: "red" | "yellow" | "blue" | "gray";
};
export function CardLogs({
  title = "Error de configuração: ",
  errorMessage,
  color = "red",
}: CardLogsProps) {
  const classColors = {
    red: "border-red-400 bg-red-50 text-red-500",
    yellow: "border-yellow-400 bg-yellow-50 text-yellow-500",
    blue: "border-blue-400 bg-blue-50 text-blue-500",
    gray: "border-gray-400 bg-gray-50 text-gray-500",
  };

  return (
    <div
      className={`mt-2 rounded-md border-[0.01rem] ${classColors[color]} pt-5 pr-4 pb-5 pl-4 text-sm font-light`}
    >
      <span className={`font-medium text-${color}-500`}>{title}</span>
      {errorMessage || "Nenhuma mensagem de erro."}
    </div>
  );
}
