import { Search } from "lucide-react";

import Input from "../Input";

export interface searchByType {
  id: string;
  name: string;
}

interface DropdownHeaderProps {
  search: string;
  setSearch: React.Dispatch<React.SetStateAction<string>>;
  setFilter: React.Dispatch<React.SetStateAction<string>>;
}

export default function DropdownHeader({
  search,
  setSearch,
  setFilter,
}: DropdownHeaderProps) {
  return (
    <div className="w-full flex-col lg:flex-row flex items-center lg:justify-between gap-6">
      <div className="w-full flex justify-center bg-primary-white">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            setFilter(search);
          }}
          className="w-full lg:w-auto lg:grow"
        >
          <Input
            type="text"
            name="search"
            placeholder="Search"
            icon={<Search />}
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
              if (e.target.value === "") {
                setFilter("");
              }
            }}
            onIconClick={() => {
              setFilter(search);
            }}
            round={true}
          />
        </form>
      </div>
    </div>
  );
}
