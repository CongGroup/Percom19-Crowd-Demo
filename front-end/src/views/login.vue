<template>
    <div class="login">
        <img class="icon" alt="Vue logo" src="../assets/cd3.jpeg"/>
        <pacman v-if="accessStatus===1"></pacman>
        <div class="form" v-else-if="accessStatus===2">
            <span class="label">Password:</span>
            <!--<input class="input" type="number" v-model.number="value">-->
            <input class="input" type="password" v-model="value">
            <button class="btn btn-dark contract-button" @click="submit"> Submit</button>
            <!--<div class="error" v-show="errors.has('value')">{{ errors.first('value') }}</div>-->
        </div>
    </div>
</template>

<script>
    // @ is an alias to /src
    const ACCESSED=0;
    const WAITING=1;
    const UNACCESSED = 2;
    require("@/assets/css/style.css");
    export default {
        name: "home",
        components: {

        },
        data: function () {
            return {
                value: undefined,
                accessStatus: UNACCESSED,
            }
        },
        methods: {
          submit:function(){
              this.accessStatus = WAITING;
              this.axios.post(`${process.env.HTTP_PATH}/validate`,this.value).then(res=>{
                  let accessToken = res.data;
                  this.$cookies.set('accessToken',accessToken,'7d');
                  this.accessStatus = ACCESSED;
                  this.$router.replace('/admin');
              }).catch(err=>{
                  console.log(err);
                  this.accessStatus = UNACCESSED;
              })
          },
        },
        computed: {

        },
        beforeCreate: function () {
        }
    };
</script>

<style scoped>
    .login {
        display: flex;
        display: -webkit-flex;
        flex-direction: column;
        align-items: center;
    }
    .login .icon {
        margin-top: 200px;
        margin-bottom: 50px;
    }
</style>
