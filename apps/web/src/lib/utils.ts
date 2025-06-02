import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const isJson = (val: string) => {
  if (!val) return true;
  try {
    JSON.parse(val);
    return true;
  } catch {
    return false;
  }
};

export const isValidXml = (val: string) => {
  if (!val.trim()) return true; // Allow empty body
  try {
    const parser = new DOMParser();
    const doc = parser.parseFromString(val, "text/xml");
    const parserError = doc.querySelector("parsererror");
    return !parserError;
  } catch {
    return false;
  }
};

export const isValidForm = (val: string): boolean => {
  const trimmed = val.trim();
  if (!trimmed) return true; // Allow empty body

  // Accept simple URL-encoded format like "key=value"
  const isURLEncoded = /^[^=&]+=[^=&]*(&[^=&]+=[^=&]*)*$/.test(trimmed);

  if (isURLEncoded) return true;

  // return isJson(trimmed);
  return false;
};
