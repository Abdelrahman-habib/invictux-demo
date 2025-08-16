"use client";

import { Input } from "@/components/ui/input";
import { Search } from "lucide-react";
import { useEffect, useState } from "react";
import { useDebounce } from "@/hooks/use-debounce";
import { cn } from "@/lib/utils";

interface SearchInputProps {
  onSearch: (query: string) => void;
  placeholder?: string;
  className?: string;
  disabled?: boolean;
  debounceTime?: number;
}

export function SearchInput({
  onSearch,
  placeholder,
  className,
  disabled,
  debounceTime = 300,
}: SearchInputProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const debouncedSearch = useDebounce(searchQuery, debounceTime);

  useEffect(() => {
    onSearch(debouncedSearch);
  }, [debouncedSearch, onSearch]);

  return (
    <div className="w-full">
      <div className="relative">
        <Search className="absolute end-3 top-3 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder={placeholder}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className={cn("ps-3", className)}
          disabled={disabled}
        />
      </div>
    </div>
  );
}
