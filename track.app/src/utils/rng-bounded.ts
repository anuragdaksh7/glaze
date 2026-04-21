export const getRNGIndexBounded = (min: number, max: number, val: string): number => {
  const hash = Array.from(val).reduce((acc, char) => acc + char.charCodeAt(0), 0);
  const range = max - min + 1;
  return (hash % range) + min;
}