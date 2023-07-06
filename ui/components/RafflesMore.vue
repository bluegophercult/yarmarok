<template>
    <HeadlessMenu as="div" class="relative inline-block text-left">
        <HeadlessMenuButton
                class="rounded-lg bg-white w-10 h-full shadow-md text-gray-600 hover:text-teal-400 transition duration-200 grid place-content-center">
            <Icon name="heroicons:ellipsis-vertical" class="w-6 h-6"/>
        </HeadlessMenuButton>

        <transition name="m-fade">
            <HeadlessMenuItems
                    class="absolute right-0 mt-2 w-48 origin-top-right divide-y divide-gray-200 rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                <div class="px-1 py-1">
                    <HeadlessMenuItem v-slot="{ active }">
                        <button :class="[
                            active ? 'bg-teal-100 text-teal-950' : 'text-gray-900',
                            'group flex w-full items-center rounded-md px-2 py-2 text-sm',
                        ]">
                            <Icon :active="active" name="heroicons:pencil" class="mr-2 h-5 w-5 text-teal-400"/>
                            Змінити
                        </button>
                    </HeadlessMenuItem>
                </div>
                <div class="px-1 py-1">
                    <HeadlessMenuItem v-slot="{ active }">
                        <button :class="[
                          active ? 'bg-teal-100 text-teal-950' : 'text-gray-900',
                          'group flex w-full items-center rounded-md px-2 py-2 text-sm',
                        ]">
                            <Icon :active="active" name="heroicons:trash" class="mr-2 h-5 w-5 text-teal-400"/>
                            Видалити
                        </button>
                    </HeadlessMenuItem>
                </div>
            </HeadlessMenuItems>
        </transition>
    </HeadlessMenu>
</template>
