import Carousel from "../../components/Carousel";
import SearchHeader from "../../components/SearchHeader";

import Categories from "./Categories";
import MostBought from "./MostBought";

import { HOMEPAGE_TITLE } from "../../const/const";
import useTitle from "../../hooks/useTitle";


const dummy = [
  "https://t4.ftcdn.net/jpg/02/92/36/71/360_F_292367179_T5xBfw6nJBwJ0HE8wfwz20QuYfOrIm8b.jpg",
  "https://www.shutterstock.com/image-vector/brush-sale-banner-vector-260nw-1090866878.jpg",
  "https://static.vecteezy.com/system/resources/previews/002/217/707/non_2x/medicine-trendy-banner-vector.jpg",
];

export default function Homepage() {
  useTitle(HOMEPAGE_TITLE)

  return (
    <>
      <SearchHeader />
      <div className="grow w-full flex flex-col bg-primary-white items-center">
        <div className="my-16 w-[90%] max-w-[1259px] flex flex-col items-center gap-10">
          <Carousel images={dummy} />
          <div className="w-full flex flex-col gap-16">
            <Categories />
            <MostBought />
          </div>
        </div>
      </div>
    </>
  );
}
