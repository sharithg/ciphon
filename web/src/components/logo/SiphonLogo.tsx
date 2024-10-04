import React from "react";

const SiphonLogo = (input: { scaleDown: number }) => {
  return (
    <svg
      width={(67 * input.scaleDown).toString()}
      height={(123 * input.scaleDown).toString()}
      viewBox="0 0 67 123"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M10.7872 119.809C53.5661 119.809 62.7479 112.491 64 63.7103V22.4681"
        stroke="white"
        stroke-width="5"
        stroke-linecap="round"
      />
      <path
        d="M49.3445 22.9622C50.6741 102.02 50.6741 116.844 10.7872 119.808"
        stroke="white"
        stroke-width="5"
        stroke-linecap="round"
      />
      <path
        d="M56.2128 3.01284C13.4339 3.01284 4.25207 10.329 3 59.1037L3 100.34"
        stroke="white"
        stroke-width="5"
        stroke-linecap="round"
      />
      <path
        d="M17.6523 99.8336C16.3227 20.7858 16.3227 5.96432 56.2096 3.00003"
        stroke="white"
        stroke-width="5"
        stroke-linecap="round"
      />
    </svg>
  );
};

export default SiphonLogo;
