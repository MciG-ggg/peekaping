import React from 'react';
import { useTimezone } from '../context/TimezoneContext';
import { Check, ChevronsUpDown } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from './ui/button';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from './ui/command';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from './ui/popover';

const fallbackTimezones = [
  'UTC',
  'Europe/London',
  'Europe/Paris',
  'America/New_York',
  'Asia/Tokyo',
  'Australia/Sydney',
  // ...add more as needed
];

function getSupportedTimezones(): string[] {
  if (
    typeof Intl !== 'undefined' &&
    'supportedValuesOf' in Intl &&
    typeof (Intl as { supportedValuesOf?: unknown }).supportedValuesOf === 'function'
  ) {
    // @ts-expect-error: supportedValuesOf is not yet in TS types
    return Intl.supportedValuesOf('timeZone');
  }
  return fallbackTimezones;
}

function getTimezoneOffsetLabel(timezone: string): string {
  try {
    const now = new Date();
    const tzDate = new Date(now.toLocaleString('en-US', { timeZone: timezone }));
    const utcDate = new Date(now.toLocaleString('en-US', { timeZone: 'UTC' }));
    const diff = (tzDate.getTime() - utcDate.getTime()) / 60000; // in minutes
    const sign = diff >= 0 ? '+' : '-';
    const absDiff = Math.abs(diff);
    const hours = Math.floor(absDiff / 60);
    const minutes = absDiff % 60;
    if (diff === 0) return 'UTC±0';
    return `UTC${sign}${hours}${minutes ? ':' + String(minutes).padStart(2, '0') : ''}`;
  } catch {
    return 'UTC';
  }
}

function getTimezoneOffsetMinutes(timezone: string): number {
  try {
    const now = new Date();
    const tzDate = new Date(now.toLocaleString('en-US', { timeZone: timezone }));
    const utcDate = new Date(now.toLocaleString('en-US', { timeZone: 'UTC' }));
    return (tzDate.getTime() - utcDate.getTime()) / 60000; // in minutes
  } catch {
    return 0;
  }
}

const timezones = getSupportedTimezones();

// Sort timezones: negative offsets first, then UTC (offset 0), then positive offsets
const sortedTimezones = [...timezones].sort((a, b) => {
  const offsetA = getTimezoneOffsetMinutes(a);
  const offsetB = getTimezoneOffsetMinutes(b);
  // Group by sign: negatives first, then zero, then positives
  if (offsetA < 0 && offsetB >= 0) return -1;
  if (offsetA >= 0 && offsetB < 0) return 1;
  if (offsetA === 0 && offsetB !== 0) return -1;
  if (offsetA !== 0 && offsetB === 0) return 1;
  return offsetA - offsetB;
});

const TimezoneSelector: React.FC = () => {
  const { timezone, setTimezone } = useTimezone();
  const [open, setOpen] = React.useState(false);
  const [search, setSearch] = React.useState('');

  // Filter by search
  const filteredTimezones = React.useMemo(() => {
    if (!search) return sortedTimezones;
    return sortedTimezones.filter(
      (tz) =>
        tz.toLowerCase().includes(search.toLowerCase()) ||
        getTimezoneOffsetLabel(tz).toLowerCase().includes(search.toLowerCase())
    );
  }, [search]);

  const selectedLabel = timezone
    ? `${timezone} (${getTimezoneOffsetLabel(timezone)})`
    : 'Select timezone...';

  return (
    <div className="flex flex-col gap-1">
      <label htmlFor="timezone-combobox" className="text-sm font-medium mb-1">
        Timezone
      </label>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="min-w-[260px] justify-between"
            id="timezone-combobox"
          >
            {selectedLabel}
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="min-w-[260px] p-0">
          <Command>
            <CommandInput
              placeholder="Search timezone..."
              value={search}
              onValueChange={setSearch}
              className="h-9"
            />
            <CommandList>
              <CommandEmpty>No timezone found.</CommandEmpty>
              <CommandGroup>
                {filteredTimezones.map((tz) => (
                  <CommandItem
                    key={tz}
                    value={tz}
                    onSelect={() => {
                      setTimezone(tz);
                      setOpen(false);
                    }}
                  >
                    <Check
                      className={cn(
                        'mr-2 h-4 w-4',
                        timezone === tz ? 'opacity-100' : 'opacity-0'
                      )}
                    />
                    {tz} <span className="ml-2 text-muted-foreground">({getTimezoneOffsetLabel(tz)})</span>
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  );
};

export default TimezoneSelector;
