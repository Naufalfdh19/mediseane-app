import { useState } from "react";
import { Search } from "lucide-react";
import { useNavigate } from "react-router-dom";

import Input from "../Input";

import { PATH_PRODUCTS } from "../../const/const";

export default function SearchHeader() {
  const [search, setSearch] = useState("");
  const navigate = useNavigate();

  return (
    <form
      className="flex justify-center p-3 md:pt-8 bg-primary-white"
      onSubmit={(e) => {
        e.preventDefault();
        const toSearch = search;
        setSearch("");
        navigate(PATH_PRODUCTS + (toSearch ? "?search=" + toSearch : ""));
      }}
    >
      <div className="w-full max-w-[794px]">
        <Input
          type="text"
          name="search"
          placeholder="Search"
          icon={<Search />}
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
          }}
          onIconClick={() => {
            const toSearch = search;
            setSearch("");
            navigate(PATH_PRODUCTS + (toSearch ? "?search=" + toSearch : ""));
          }}
          round={true}
        />
      </div>
    </form>
  );
}
