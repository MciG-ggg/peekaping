import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { X } from "lucide-react";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Checkbox } from "@/components/ui/checkbox";
import { useFormContext } from "react-hook-form";
import { z } from "zod";
import { useQuery } from "@tanstack/react-query";
import { getTagsOptions } from "@/api/@tanstack/react-query.gen";
import { type TagModel } from "@/api";
import { useState } from "react";

export const tagsDefaultValues = {
  tag_ids: [],
};

export const tagsSchema = z.object({
  tag_ids: z.array(z.string()).default([]),
});

const Tags = () => {
  const form = useFormContext();
  const [tagPopoverOpen, setTagPopoverOpen] = useState(false);

  // Load available tags
  const { data: tagsData } = useQuery({
    ...getTagsOptions({
      query: {
        limit: 100,
      },
    }),
  });

  const availableTags = (tagsData?.data || []) as TagModel[];
  const selectedTagIds = form.watch("tag_ids") || [];
  const selectedTags = availableTags.filter((tag) =>
    selectedTagIds.includes(tag.id!)
  );

  const handleTagToggle = (tagId: string) => {
    const currentTags = form.getValues("tag_ids") || [];
    const newTags = currentTags.includes(tagId)
      ? currentTags.filter((id: string) => id !== tagId)
      : [...currentTags, tagId];
    form.setValue("tag_ids", newTags);
  };

  const handleTagRemove = (tagId: string) => {
    const currentTags = form.getValues("tag_ids") || [];
    const newTags = currentTags.filter((id: string) => id !== tagId);
    form.setValue("tag_ids", newTags);
  };

  const clearAllTags = () => {
    form.setValue("tag_ids", []);
  };

  return (
    <FormField
      control={form.control}
      name="tag_ids"
      render={() => (
        <FormItem>
          <FormLabel>Tags</FormLabel>
          <FormControl>
            <div className="space-y-2">
              <Popover open={tagPopoverOpen} onOpenChange={setTagPopoverOpen}>
                <PopoverTrigger asChild>
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full justify-start text-left font-normal"
                  >
                    {selectedTags.length > 0 ? (
                      <div className="flex flex-wrap gap-1">
                        {selectedTags.slice(0, 3).map((tag) => (
                          <Badge
                            key={tag.id}
                            variant="secondary"
                            className="text-xs"
                            style={{ backgroundColor: tag.color, color: 'white' }}
                          >
                            {tag.name}
                          </Badge>
                        ))}
                        {selectedTags.length > 3 && (
                          <Badge variant="secondary" className="text-xs">
                            +{selectedTags.length - 3} more
                          </Badge>
                        )}
                      </div>
                    ) : (
                      <span className="text-muted-foreground">Select tags...</span>
                    )}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-80 p-0">
                  <div className="max-h-60 overflow-y-auto">
                    <div className="p-2">
                      {availableTags.map((tag) => (
                        <div
                          key={tag.id}
                          className="flex items-center space-x-2 p-2 hover:bg-accent hover:text-accent-foreground rounded-sm cursor-pointer"
                          onClick={() => handleTagToggle(tag.id!)}
                        >
                          <Checkbox
                            checked={selectedTagIds.includes(tag.id!)}
                            onChange={() => handleTagToggle(tag.id!)}
                          />
                          <Badge
                            variant="secondary"
                            className="text-xs"
                            style={{ backgroundColor: tag.color, color: 'white' }}
                          >
                            {tag.name}
                          </Badge>
                          {/* {tag.description && (
                            <span className="text-xs text-muted-foreground">
                              {tag.description}
                            </span>
                          )} */}
                        </div>
                      ))}
                      {availableTags.length === 0 && (
                        <div className="text-center text-muted-foreground text-sm py-4">
                          No tags available
                        </div>
                      )}
                    </div>
                  </div>
                </PopoverContent>
              </Popover>

              {/* Show selected tags below the button */}
              {selectedTags.length > 0 && (
                <div className="flex flex-wrap gap-1">
                  {selectedTags.map((tag) => (
                    <Badge
                      key={tag.id}
                      variant="secondary"
                      className="text-xs flex items-center gap-1"
                      style={{ backgroundColor: tag.color, color: 'white' }}
                    >
                      {tag.name}
                      <X
                        className="h-3 w-3 cursor-pointer"
                        onClick={() => handleTagRemove(tag.id!)}
                      />
                    </Badge>
                  ))}
                </div>
              )}
            </div>
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};

export default Tags;
