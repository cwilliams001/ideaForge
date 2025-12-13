"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface DeleteConfirmModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: () => void;
  title: string;
  category: string;
  isDeleting: boolean;
}

export function DeleteConfirmModal({
  open,
  onOpenChange,
  onConfirm,
  title,
  category,
  isDeleting,
}: DeleteConfirmModalProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md border-2 border-destructive/50">
        <DialogHeader>
          <DialogTitle className="font-mono uppercase text-destructive">
            [ DELETE CONFIRMATION ]
          </DialogTitle>
          <DialogDescription className="font-mono text-sm space-y-2 pt-4">
            <p className="text-foreground">Are you sure you want to delete this note?</p>
            <div className="bg-destructive/10 border border-destructive/30 p-3 rounded space-y-1">
              <p><span className="text-muted-foreground">Title:</span> {title}</p>
              <p><span className="text-muted-foreground">Category:</span> {category}</p>
            </div>
            <p className="text-destructive text-xs pt-2">
              ⚠ This action cannot be undone.
            </p>
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="flex-row gap-2 sm:justify-end">
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isDeleting}
            className="uppercase"
          >
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={onConfirm}
            disabled={isDeleting}
            className="uppercase"
          >
            {isDeleting ? "Deleting..." : "[ Delete ]"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
