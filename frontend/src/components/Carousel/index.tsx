import { useState, useEffect } from "react";

interface CarouselProps {
  images: string[];
  autoSlide?: boolean;
  autoSlideInterval?: number;
}

export default function Carousel({
  images,
  autoSlide = true,
  autoSlideInterval = 3000,
}: CarouselProps) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [resetTimer, setResetTimer] = useState(false);

  useEffect(() => {
    if (autoSlide) {
      const slideInterval = setInterval(() => {
        setCurrentIndex((prevIndex) => (prevIndex + 1) % images.length);
      }, autoSlideInterval);
      return () => clearInterval(slideInterval);
    }
  }, [autoSlide, autoSlideInterval, images.length, resetTimer]);

  const nextSlide = () => {
    setCurrentIndex((prevIndex) => (prevIndex + 1) % images.length);
  };

  const prevSlide = () => {
    setCurrentIndex(
      (prevIndex) => (prevIndex - 1 + images.length) % images.length
    );
  };

  return (
    <div className="relative w-full mx-auto">
      <div className="overflow-hidden relative h-32 md:h-64">
        {images.map((image, index) => (
          <div
            key={index}
            className={`absolute inset-0 transition-transform transform ${
              index === currentIndex ? "translate-x-0" : "translate-x-full"
            }`}
          >
            <img
              src={image}
              alt={`Slide ${index}`}
              className="w-full h-full object-cover"
            />
          </div>
        ))}
      </div>
      <button
        aria-label="Previous page"
        className="absolute top-1/2 left-0 transform -translate-y-1/2 bg-gray-800 text-primary-white p-2 h-full w-[25%]"
        onClick={() => {
          setResetTimer((prev) => !prev);
          prevSlide();
        }}
      ></button>
      <button
        aria-label="Next page"
        className="absolute top-1/2 right-0 transform -translate-y-1/2 bg-gray-800 text-primary-white p-2 h-full w-[25%]"
        onClick={() => {
          setResetTimer((prev) => !prev);
          nextSlide();
        }}
      ></button>
      <div className="absolute bottom-0 left-0 right-0 flex justify-center mb-4 opacity-80">
        {images.map((_, index) => (
          <div
            key={index}
            className={`w-2 h-2 rounded-full mx-1 ${
              index === currentIndex ? "bg-primary-gray" : "bg-secondary-gray"
            }`}
            onClick={() => {
              setResetTimer((prev) => !prev);
              setCurrentIndex(index);
            }}
          />
        ))}
      </div>
    </div>
  );
}
