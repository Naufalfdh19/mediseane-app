import {
  ArrowDownNarrowWide,
  ArrowUpNarrowWide,
  Plus,
  Search,
} from "lucide-react";

import Button from "../Button";
import Input from "../Input";
import Dropdown from "../Dropdown";
import { CapitalizeFirstLetter } from "../../utils/formatter";

export interface searchByType {
  id: string;
  name: string;
}

export interface sortByType {
  id: string;
  name: string;
}

interface TableHeaderProps {
  search: string;
  searchBy?: searchByType;
  setSearchBy?: React.Dispatch<React.SetStateAction<searchByType>>;
  searchByList?: searchByType[];
  setSearch: React.Dispatch<React.SetStateAction<string>>;
  setFilter: React.Dispatch<React.SetStateAction<string>>;
  sortBy: sortByType | undefined;
  setSortBy: React.Dispatch<React.SetStateAction<sortByType | undefined>>;
  sortByList: sortByType[];
  asc: boolean;
  setAsc: React.Dispatch<React.SetStateAction<boolean>>;
  onAdd?: () => void;
  readOnly?: boolean;
}

export default function TableHeader({
  search,
  searchBy = { id: "", name: "" },
  setSearchBy,
  searchByList = [],
  setSearch,
  setFilter,
  sortBy,
  setSortBy,
  sortByList,
  asc,
  setAsc,
  onAdd = () => {},
  readOnly = false,
}: TableHeaderProps) {
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
      {searchByList && searchByList.length > 0 ? (
        <div className="w-full flex-col sm:flex-row flex items-center lg:justify-between gap-6">
          <div className="w-full lg:max-w-[260px] flex gap-2.5 items-center">
            <div className="w-full">
              <Dropdown
                name="search-by"
                placeholder="Search By"
                data={searchByList.map((item) => {
                  return {
                    id: item.id,
                    name: CapitalizeFirstLetter(item.name.split("_").join(" ")),
                  };
                })}
                selectedId={searchBy.id}
                onSelect={(id) => {
                  const searching = searchByList.find((item) => item.id === id);
                  if (searching) {
                    setSearchBy && setSearchBy(searching);
                  }
                }}
              />
            </div>
          </div>
          <div className="w-full lg:max-w-[260px] flex gap-2.5 items-center">
            <div className="w-[48px]">
              <Button
                size="sm"
                type="ghost-green"
                square={true}
                onClick={() => {
                  setAsc((prev) => !prev);
                }}
              >
                <div className="flex justify-center">
                  {asc ? (
                    <ArrowUpNarrowWide aria-label="Ascending" />
                  ) : (
                    <ArrowDownNarrowWide aria-label="Descending" />
                  )}
                </div>
              </Button>
            </div>
            <div className="w-full">
              <Dropdown
                name="sort-by"
                placeholder="Sort By"
                data={sortByList.map((item) => {
                  return {
                    id: item.id,
                    name: CapitalizeFirstLetter(item.name.split("_").join(" ")),
                  };
                })}
                selectedId={sortBy ? sortBy.id : undefined}
                onSelect={(id) => {
                  const sorting = sortByList.find((item) => item.id === id);
                  if (sorting) {
                    setSortBy(sorting);
                  }
                }}
              />
            </div>
          </div>
        </div>
      ) : (
        <div className="w-full lg:max-w-[260px] flex gap-2.5 items-center">
          <div className="w-[48px]">
            <Button
              size="sm"
              type="ghost-green"
              square={true}
              onClick={() => {
                setAsc((prev) => !prev);
              }}
            >
              <div className="flex justify-center">
                {asc ? <ArrowUpNarrowWide /> : <ArrowDownNarrowWide />}
              </div>
            </Button>
          </div>
          <div className="w-full">
            <Dropdown
              name="sort-by"
              placeholder="Sort By"
              data={sortByList.map((item) => {
                return {
                  id: item.id,
                  name: CapitalizeFirstLetter(item.name.split("_").join(" ")),
                };
              })}
              selectedId={sortBy ? sortBy.id : undefined}
              onSelect={(id) => {
                const sorting = sortByList.find((item) => item.id === id);
                if (sorting) {
                  setSortBy(sorting);
                }
              }}
            />
          </div>
        </div>
      )}
      {!readOnly && (
        <div className="w-full lg:max-w-[125px]">
          <Button onClick={onAdd}>
            <div className="flex gap-1.5 justify-center px-2">
              <Plus />
              <h1>Add</h1>
            </div>
          </Button>
        </div>
      )}
    </div>
  );
}
