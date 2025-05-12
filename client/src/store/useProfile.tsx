import { create } from "zustand";

type Profile = {
  username: string;
  email: string;
};

type ProfileStore = Profile & {
  setProfile: (profile: Profile) => void;
  clearProfile: () => void;
};

export const useProfile = create<ProfileStore>((set) => ({
  username: "",
  email: "",
  setProfile: (profile) =>
    set(() => ({
      username: profile.username,
      email: profile.email,
    })),
  clearProfile: () => set(() => ({ username: "", email: "" })),
}));
