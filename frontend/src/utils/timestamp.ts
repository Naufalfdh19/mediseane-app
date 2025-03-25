export function parseTimestamp(create_timestamp: string) {
  const createdAt = new Date(create_timestamp);
  createdAt.setHours(createdAt.getHours() - 7);
  const createdDate = createdAt.toLocaleDateString("en-US");
  const createdTime = createdAt.toLocaleTimeString("en-US");

  return `${createdDate} at ${createdTime}`;
}
