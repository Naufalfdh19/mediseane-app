export default function LoaderCard() {
  return (
    <div className="w-[90%] lg:w-[80%] mx-auto flex justify-center items-center gap-32">
      <div className="flex flex-row gap-2">
        <div className="w-4 h-4 rounded-full bg-primary-green animate-bounce"></div>
        <div className="w-4 h-4 rounded-full bg-primary-green animate-bounce [animation-delay:-.3s]"></div>
        <div className="w-4 h-4 rounded-full bg-primary-green animate-bounce [animation-delay:-.5s]"></div>
      </div>
    </div>
  );
}
