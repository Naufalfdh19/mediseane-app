import { useEffect, useRef, useState } from "react";
import useOutsideClick from "../../hooks/useOutsideClick";
import { ChevronDown } from "lucide-react";

interface DropdownItem {
  id: string;
  name: string;
  imageUrl?: string;
}

interface DropdownProps {
  name: string;
  placeholder?: string;
  data: DropdownItem[];
  hasImage?: boolean;
  style?: string;
  selectedId?: string;
  onSelect?: (id: string) => void;
  valid?: boolean;
  disabled?: boolean;
}

const Dropdown = ({
  name,
  placeholder = "Select",
  data,
  hasImage,
  style,
  selectedId,
  onSelect,
  valid = true,
  disabled = false,
}: DropdownProps) => {
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [selectedItem, setSelectedItem] = useState<DropdownItem | undefined>(
    selectedId ? data?.find((item) => item.id === selectedId) : undefined
  );

  const handleChange = (item: DropdownItem) => {
    setSelectedItem(item);
    onSelect && onSelect(item.id);
    setIsOpen(false);
  };

  useEffect(() => {
    if (selectedId && data) {
      const newSelectedItem = data.find((item) => item.id === selectedId);
      newSelectedItem && setSelectedItem(newSelectedItem);
    } else {
      setSelectedItem(undefined);
    }
  }, [selectedId, data]);

  const dropdownRef = useRef<HTMLDivElement>(null);
  useOutsideClick({
    ref: dropdownRef,
    handler: () => setIsOpen(false),
  });

  const dropdownClass =
    "absolute bg-primary-white w-full max-h-52 overflow-y-auto py-3 rounded shadow-md z-10 top-full left-0 mt-1";

  return (
    <div
      ref={dropdownRef}
      className={`relative w-full flex border px-[13px] py-[12px] rounded-lg gap-2 ${
        valid ? "border-primary-black" : "border-primary-red"
      } ${disabled ? "bg-secondary-gray" : ""}`}
    >
      <div className="w-full">
        <button
          id={name}
          aria-label="Toggle dropdown"
          aria-haspopup="true"
          aria-expanded={isOpen && !disabled}
          type="button"
          onClick={() => setIsOpen(!isOpen)}
          disabled={disabled}
          className={`flex justify-between items-center gap-5 rounded w-full ${style}`}
        >
          <span
            className={`truncate
              ${
                selectedItem?.name ? "text-primary-black" : "text-primary-gray"
              }`}
          >
            {selectedItem?.name || placeholder}
          </span>
          <ChevronDown
            size={20}
            className={`transform duration-500 ease-in-out ${
              isOpen && !disabled ? "rotate-180" : ""
            }`}
          />
        </button>
        {/* Open */}
      </div>
      {isOpen && !disabled && (
        <div aria-label="Dropdown menu" className={dropdownClass}>
          <ul
            role="menu"
            aria-labelledby={name}
            aria-orientation="vertical"
            className="leading-10"
          >
            {data?.map((item) => (
              <li
                key={item.id}
                onClick={() => handleChange(item)}
                className={`flex items-center cursor-pointer hover:bg-primary-gray hover:text-primary-white px-3 ${
                  selectedItem?.id === item.id
                    ? "bg-primary-blue text-primary-white"
                    : ""
                }`}
              >
                {hasImage && (
                  <img
                    src={item.imageUrl}
                    alt="image"
                    loading="lazy"
                    className="w-8 h-8 rounded-full bg-primary-gray object-cover me-2"
                  />
                )}
                <span>{item.name}</span>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default Dropdown;
