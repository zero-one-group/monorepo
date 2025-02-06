// SVG by Sam Herbert (@https://github.com/SamHerbert/SVG-Loaders)
// Alternatively you can use https://uiball.com/ldrs/
const SVGLoader = () => (
  <svg
    width={40}
    height={30}
    viewBox="0 0 120 30"
    xmlns="http://www.w3.org/2000/svg"
    fill="#020617"
  >
    <circle cx={15} cy={15} r={15}>
      <animate
        attributeName="r"
        from={15}
        to={15}
        begin="0s"
        dur="0.8s"
        values="15;9;15"
        calcMode="linear"
        repeatCount="indefinite"
      />
      <animate
        attributeName="fill-opacity"
        from={1}
        to={1}
        begin="0s"
        dur="0.8s"
        values="1;.5;1"
        calcMode="linear"
        repeatCount="indefinite"
      />
    </circle>
    <circle cx={60} cy={15} r={9} fillOpacity="0.3">
      <animate
        attributeName="r"
        from={9}
        to={9}
        begin="0s"
        dur="0.8s"
        values="9;15;9"
        calcMode="linear"
        repeatCount="indefinite"
      />
      <animate
        attributeName="fill-opacity"
        from="0.5"
        to="0.5"
        begin="0s"
        dur="0.8s"
        values=".5;1;.5"
        calcMode="linear"
        repeatCount="indefinite"
      />
    </circle>
    <circle cx={105} cy={15} r={15}>
      <animate
        attributeName="r"
        from={15}
        to={15}
        begin="0s"
        dur="0.8s"
        values="15;9;15"
        calcMode="linear"
        repeatCount="indefinite"
      />
      <animate
        attributeName="fill-opacity"
        from={1}
        to={1}
        begin="0s"
        dur="0.8s"
        values="1;.5;1"
        calcMode="linear"
        repeatCount="indefinite"
      />
    </circle>
  </svg>
)

export default function PageLoader() {
  return (
    <div className="flex size-full min-h-screen flex-col items-center justify-center p-4 py-12 sm:px-6 lg:px-8">
      <div className="flex flex-col items-center justify-center space-y-4 sm:mx-auto sm:w-full sm:max-w-lg">
        <h1 className="mt-3 text-center font-medium">Loading...</h1>
        <p className="text-center text-neutral-700 text-sm leading-6 tracking-tight">
          Does this take longer than expected? <br />
          Try clearing your browser's cache or check if you have an ad blocker enabled!
        </p>
        <div className="mt-8">
          <SVGLoader />
        </div>
      </div>
    </div>
  )
}
