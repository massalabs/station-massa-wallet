export function Loading() {
  return (
    <div
      className="flex flex-col justify-center items-center gap-5 w-1/2"
      data-testid="loading"
    >
      <div className="bg-secondary rounded-2xl w-full max-w-lg p-12">
        <div className="h-2 w-20 bg-c-disabled-1 rounded-lg animate-pulse mb-4"></div>
        <div data-testid="balance" className="flex items-center w-fit mb-6">
          <div className="h-10 w-10 bg-[red] rounded-full animate-pulse mr-2 blur-sm"></div>
          <label className="mas-banner mb-2 text-f-primary blur-md animate-pulse">
            000,000.00
          </label>
        </div>
        <div className="flex gap-7">
          <div className="h-12 w-48 bg-c-disabled-1 rounded-lg animate-pulse"></div>
          <div className="h-12 w-48 bg-c-disabled-1 rounded-lg animate-pulse"></div>
        </div>
      </div>
      <div className="bg-secondary rounded-2xl w-full max-w-lg p-10">
        <div className="h-2 w-20 bg-c-disabled-1 rounded-lg animate-pulse mb-6"></div>
        <div
          data-testid="clipboard-field"
          className="flex flex-row items-center mas-body2 justify-between
              w-full h-12 px-3 rounded bg-primary cursor-pointer"
        >
          <div className="h-2 w-16 bg-c-disabled-1 rounded-lg animate-pulse"></div>
        </div>
      </div>
    </div>
  );
}
