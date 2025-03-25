import { ChevronLeft, ChevronRight } from "lucide-react";
import { SetStateAction } from "react";

import Button from "../Button";

interface PaginationProps {
  page: number;
  setPage: React.Dispatch<SetStateAction<number>>;
  totalPage: number;
}

export default function Pagination(props: PaginationProps) {
  function range(start: number, end: number) {
    return Array.from({ length: end - start + 1 }, (_, i) => i + start);
  }

  return (
    <div className="flex flex-wrap justify-center gap-3 items-center">
      <ChevronLeft
        className="cursor-pointer"
        onClick={() => {
          props.setPage((page) => (page > 1 ? page - 1 : page));
        }}
      />
      {range(1, props.totalPage)
        .filter((num) => {
          if (props.page === 1) {
            return (
              num === 1 || num <= props.page + 3 || num === props.totalPage
            );
          }

          if (props.page === props.totalPage) {
            return (
              num === 1 || num >= props.page - 3 || num === props.totalPage
            );
          }

          return (
            num === 1 ||
            (num >= props.page - 2 && num <= props.page + 2) ||
            num === props.totalPage
          );
        })
        .map((num) => {
          if (num !== 1 && num !== props.totalPage) {
            if (props.page === 1 && num === props.page + 3) {
              return (
                <p key={num} className="text-primary-black">
                  ...
                </p>
              );
            } else if (
              props.page === props.totalPage &&
              num === props.page - 3
            ) {
              return (
                <p key={num} className="text-primary-black">
                  ...
                </p>
              );
            }
            if (
              !(props.page === 1 || props.page === props.totalPage) &&
              (num === props.page - 2 || num === props.page + 2)
            ) {
              return (
                <p key={num} className="text-primary-black">
                  ...
                </p>
              );
            }
          }

          return (
            <div key={num} className="min-w-[30px]">
              <Button
                size="xs"
                square={true}
                type={props.page === num ? "default" : "ghost-green"}
                submit={false}
                onClick={() => {
                  props.setPage(num);
                }}
              >
                <h1 className="px-2">{num}</h1>
              </Button>
            </div>
          );
        })}
      <ChevronRight
        className="cursor-pointer"
        onClick={() => {
          props.setPage((page) => (page < props.totalPage ? page + 1 : page));
        }}
      />
    </div>
  );
}
