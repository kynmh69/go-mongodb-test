'use client';

import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { userApi, ApiError } from '@/lib/api';
import { User, UpdateUserRequest } from '@/lib/types';

interface UpdateUserDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  user: User;
  onSuccess: () => void;
}

export function UpdateUserDialog({ open, onOpenChange, user, onSuccess }: UpdateUserDialogProps) {
  const [formData, setFormData] = useState({
    user_id: '',
    email: '',
    password: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (open && user) {
      setFormData({
        user_id: user.user_id,
        email: user.email,
        password: '',
      });
      setError(null);
    }
  }, [open, user]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const updateData: UpdateUserRequest = {};
    
    if (formData.user_id !== user.user_id) {
      updateData.user_id = formData.user_id;
    }
    
    if (formData.email !== user.email) {
      updateData.email = formData.email;
    }
    
    if (formData.password) {
      if (formData.password.length < 6) {
        setError('Password must be at least 6 characters');
        return;
      }
      updateData.password = formData.password;
    }

    if (Object.keys(updateData).length === 0) {
      setError('No changes detected');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      await userApi.updateUser(user.id, updateData);
      handleClose();
      onSuccess();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Failed to update user');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setError(null);
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Update User</DialogTitle>
          <DialogDescription>
            Update user information. Leave password blank to keep current password.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="update_user_id">User ID</Label>
              <Input
                id="update_user_id"
                value={formData.user_id}
                onChange={(e) => setFormData({ ...formData, user_id: e.target.value })}
                placeholder="Enter user ID"
                required
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="update_email">Email</Label>
              <Input
                id="update_email"
                type="email"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                placeholder="Enter email address"
                required
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="update_password">New Password (optional)</Label>
              <Input
                id="update_password"
                type="password"
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                placeholder="Enter new password (leave blank to keep current)"
              />
            </div>
            {error && (
              <div className="text-sm text-destructive bg-destructive/10 p-2 rounded">
                {error}
              </div>
            )}
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={handleClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? 'Updating...' : 'Update User'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}