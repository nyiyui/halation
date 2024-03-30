import { writable } from 'svelte/store'

export function storable(name, data) {
   const store = writable(data);
   const { subscribe, set, update } = store;
   const isBrowser = typeof window !== 'undefined';

   isBrowser &&
      localStorage.getItem(name) != null &&
      set(JSON.parse(localStorage.getItem(name)));

   return {
      subscribe,
      set: n => {
         isBrowser && (localStorage.setItem(name, JSON.stringify(n)));
         set(n);
      },
      update: cb => {
         const updatedStore = cb(get(store));

         isBrowser && (localStorage.setItem(name, JSON.stringify(updatedStore)));
         set(updatedStore);
      }
   };
}
