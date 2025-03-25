import { CapitalizeFirstLetter } from "../../utils/formatter";
import { Err } from "../../types/response";

interface ErrorCardProps {
  errors: Err[];
  noMargin?: boolean;
}

export default function ErrorCard({
  errors,
  noMargin = false,
}: ErrorCardProps) {
  return (
    <div
      className={`w-full border border-primary-red rounded-lg px-[13px] py-[12px] bg-secondary-red ${
        !noMargin && "mb-8"
      }`}
    >
      {errors.map((error, i) => {
        return (
          <p key={i} className="text-primary-red">
            {error.field === "server"
              ? "Sorry, something went wrong please try again or reload the page"
              : CapitalizeFirstLetter(error.field) +
                ": " +
                CapitalizeFirstLetter(error.detail)}
          </p>
        );
      })}
    </div>
  );
}
