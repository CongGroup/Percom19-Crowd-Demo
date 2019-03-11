<template>
    <div class="preLoading">
        <div class="loading">
            <div>
                Initializing wallet with test ether
            </div>
            <div class="pacman">
                <pacman> </pacman>
            </div>
        </div>
    </div>
</template>

<script>
    // @ is an alias to /src
    require("@/assets/css/style.css");
    export default {
        name: "walletLoading",
        data: function () {
            return {
            }
        },
    };
</script>
