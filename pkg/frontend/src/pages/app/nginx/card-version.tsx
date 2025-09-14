type CardVersionProps = {
  name: string;
  version: string;
};
export function CardVersion({ name, version }: CardVersionProps) {
  return (
    <div className="rounded-lg bg-white p-5">
      <p className="mb-1 text-base font-medium">{name}</p>
      <p className="text-[0.8rem] font-light text-gray-600">{version}</p>
    </div>
  );
}
