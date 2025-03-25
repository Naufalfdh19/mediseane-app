import { CircleAlert } from "lucide-react";
import Button from "../../components/Button";

export default function ServerError() {
  function handleTryAgain() {
    window.location.reload();
  }

  return (
    <div className="h-screen flex flex-col">
      <div className="grow flex flex-col justify-center items-center bg-primary-white">
        <div className="w-[90%] lg:w-[80%] mx-auto flex flex-col items-center gap-3">
          <CircleAlert
            className="h-32 w-32 md:h-64 md:w-64 text-primary-red"
            strokeWidth={0.5}
          />
          <h1 className="text-center text-2xl md:text-4xl font-bold text-primary-black">
            Oops...
          </h1>
          <p className="text-center md:text-xl text-primary-black">
            Sorry, something went wrong
          </p>
          <div className="w-[80%] md:w-[500px] md:mt-5">
            <Button onClick={handleTryAgain}>Try again</Button>
          </div>
        </div>
      </div>
    </div>
  );
}
