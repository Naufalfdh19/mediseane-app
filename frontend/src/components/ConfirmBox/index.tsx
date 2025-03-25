import { CircleHelp, TriangleAlert } from "lucide-react";

import Modal from "../Modal/Modal";
import Button from "../Button";

interface ConfirmBoxProps {
  onModalClick?: () => void;
  type?:
    | "delete"
    | "save"
    | "update"
    | "logout"
    | "ship"
    | "confirm"
    | "cancel";
  onYes?: () => void;
  onCancel?: () => void;
}

export default function ConfirmBox({
  type = "save",
  onYes = () => {},
  onCancel = () => {},
}: ConfirmBoxProps) {
  return (
    <Modal onClose={onCancel}>
      <div className="w-full py-8 px-2 flex flex-col items-center justify-center gap-6">
        {type === "delete" || type === "cancel" ? (
          <TriangleAlert className="w-[100px] h-[100px] md:w-[130px] md:h-[130px] text-primary-red" />
        ) : (
          <CircleHelp className="w-[100px] h-[100px] md:w-[130px] md:h-[130px] text-primary-green" />
        )}
        <h1 className="px-2 font-semibold text-primary-black text-2xl md:text-3xl text-center">
          {type === "logout"
            ? "Do you want to logout?"
            : type === "delete"
            ? "Delete entry?"
            : type === "confirm"
            ? "Confirm this order?"
            : type === "cancel"
            ? "Cancel this order?"
            : type === "ship"
            ? "Ship this order?"
            : type === "update"
            ? "Save changes?"
            : "Save new entry?"}
        </h1>
        <div className="w-full flex flex-wrap gap-3 md:gap-7 justify-center items-center">
          <div className="w-[120px] md:w-[163px]">
            <Button
              submit={false}
              type={
                type === "delete" || type === "cancel"
                  ? "default-red"
                  : "default"
              }
              size="md"
              onClick={onYes}
            >
              Yes
            </Button>
          </div>
          <div className="w-[120px] md:w-[163px]">
            <Button submit={false} type="ghost" size="md" onClick={onCancel}>
              Cancel
            </Button>
          </div>
        </div>
      </div>
    </Modal>
  );
}
