import { useEffect, useRef } from "react";

const useOutsideClickHandler = (handler: (...args: any) => any) => {
  const domNodeRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      // @ts-ignore
      if (domNodeRef.current && !domNodeRef?.current?.contains(event.target)) {
        handler();
      }
    };

    document.addEventListener("mousedown", handleClickOutside);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [handler]);

  return domNodeRef;
};

export default useOutsideClickHandler;
