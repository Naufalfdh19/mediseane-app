export function IsHourRangeValid(start: string, end: string) {
  const endSplit = end.split(":");
  const startSplit = start.split(":");
  if (endSplit[0] === startSplit[0]) {
    return endSplit[1] > startSplit[1];
  }
  return endSplit[0] > startSplit[0];
}
