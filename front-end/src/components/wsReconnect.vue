<template>
    <div class="preLoading">
        <div>
            Network down. Reconnecting
        </div>
        <div class="pacman">
            <pacman> </pacman>
        </div>
    </div>
</template>

<script>
    // @ is an alias to /src
    require("@/assets/css/style.css");
    export default {
        name: "wsReconnect",
        data: function () {
            return {
            }
        },
    };
</script>