type CardStatusProps = {
  name: string;
  status: boolean | null;
  statusMessage: string;
};
export function CardStatus({
  name,
  status,
  statusMessage: message,
}: CardStatusProps) {
  const bgColor = {
    false: "bg-red-500",
    true: "bg-green-500",
    wait: "bg-yellow-500",
  }[status === null ? "wait" : status.toString()];

  return (
    <div className="rounded-lg bg-white p-5">
      <p className="mb-1 text-base font-medium">
        <span
          className={`mr-2 inline-block h-3 w-3 rounded-full ${bgColor}`}
        ></span>
        {name}
      </p>
      <p className="text-[0.8rem] font-light text-gray-600">
        {status !== null && message}
      </p>
    </div>
  );
}
