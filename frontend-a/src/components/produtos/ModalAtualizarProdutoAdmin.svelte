<script lang="ts">
  import type {
    Categoria,
    Produto,
    ProdutoUnico,
  } from "../../lib/types/produtos";
  import { url } from "../../lib/url";
  import {
    Button,
    Modal,
    Label,
    Input,
    Checkbox,
    MultiSelect,
  } from "flowbite-svelte";
  import toast from "svelte-french-toast";

  export let formModal: boolean;
  export let IdProduto: string | undefined;
  export let categorias: Categoria[] | undefined;

  let categoriaMultiSelect: { value: string; name: string }[] = [];
  let selectedCategory: string[];

  if (categorias) {
    for (let i = 0; i < categorias.length; i++) {
      categoriaMultiSelect.push({
        value: categorias[i].Nome,
        name: categorias[i].Nome,
      });
    }
  } else {
    toast.error("Erro ao carregar categorias");
  }

  const pegarProdutos = async () => {
    return fetch(`${url.produtos}/${IdProduto}`, {
      method: "GET",
      // mode: 'no-cors',
    })
      .then(async (res) => {
        if (res.status === 200) {
          const body = await res.json();
          const produtos = body as ProdutoUnico;
          selectedCategory = produtos.CategoriasRelacionadas.split(",");

          return produtos;
        }
        toast.error(await res.text());
      })
      .catch(async (err) => {
        toast.error(err);
        console.log(err);
      });
  };
</script>

<Modal
  bind:open={formModal}
  size="lg"
  autoclose={false}
  title="Atualizar Produto"
  class="w-full shadow-lg shadow-slate-950"
>
  <form id="atualizarProduto" class="flex flex-col space-y-6">
    <!-- Modal body -->
    <!-- <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Atualizar Produto</h3> -->

    {#await pegarProdutos()}
      <div class="grid gap-4 mb-4 sm:grid-cols-3 animate-pulse">
        <div class="col-span-3">
          <div class="bg-gray-700 h-6 rounded-lg" />
        </div>
        <div>
          <div class="bg-gray-700 h-6 rounded-lg" />
        </div>
        <div>
          <div class="bg-gray-700 h-6 rounded-lg" />
        </div>
        <div class="col-span-1">
          <div class="bg-gray-700 h-6 rounded-lg" />
        </div>
        <div>
          <div class="bg-gray-700 h-6 rounded-lg" />
        </div>
        <div class="col-span-2" />
        <div class="col-span-3">
          <div class="bg-gray-700 h-6 rounded-lg" />
        </div>
        <div class="sm:col-span-3">
          <div class="bg-gray-700 h-24 rounded-lg" />
        </div>
        <div class="flex justify-center col-span-1">
          <div class="bg-gray-700 h-24 w-24 rounded-lg" />
        </div>
        <div class="col-span-2">
          <div class="bg-gray-700 h-6 rounded-lg" />
        </div>
        <div
          class="flex items-center pl-4 border border-gray-200 rounded dark:border-gray-700"
        >
          <div class="bg-gray-700 h-4 w-4 my-5 rounded-full" />
          <div class="bg-gray-700 h-6 w-20 ml-2 rounded-lg" />
        </div>
      </div>
    {:then res}
      {#if res}
        <div class="grid gap-4 mb-4 sm:grid-cols-3">
          <div class="col-span-3">
            <label
              for="nome"
              class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >Name</label
            >
            <input
              value={res.ProdutoNome}
              type="text"
              name="nome"
              id="nome"
              class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"
              placeholder="Nome do produto"
              required
              autocomplete="off"
            />
          </div>
          <div>
            <label
              for="unidade-medida"
              class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >Unidade de medida</label
            >
            <select
              value={res.ProdutoUnidadeMedida}
              id="unidade-medida"
              required
              name="unidade-medida"
              class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"
            >
              <option value="G">G</option>
              <option value="KG">KG</option>
            </select>
          </div>
          <div>
            <label
              for="qt-pacote"
              class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >Quantidade por pacote</label
            >
            <input
              value={res.ProdutoQuantidadePacote}
              type="number"
              name="qt-pacote"
              id="qt-pacote"
              class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"
              placeholder="Nome do produto"
              required
              autocomplete="off"
            />
          </div>
          <div class="col-span-1">
            <label
              for="peso"
              class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >Peso</label
            >
            <input
              value={res.ProdutoPeso}
              type="number"
              name="peso"
              id="peso"
              class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"
              placeholder="Nome do produto"
              required
              autocomplete="off"
            />
          </div>
          <div>
            <label
              for="preco"
              class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >Preço</label
            >
            <input
              value={res.ProdutoPreco}
              type="number"
              name="preco"
              id="preco"
              class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"
              placeholder="1.99"
              required
              autocomplete="off"
            />
          </div>
          <div class="col-span-2" />

          <div class="col-span-3">
            <label
              for="categoria"
              class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >Categorias</label
            >
            <MultiSelect
              id="categoria"
              name="categoria"
              items={categoriaMultiSelect}
              bind:value={selectedCategory}
              size="md"
              required
            />
          </div>
          <div class="sm:col-span-3">
            <label
              for="descricao"
              class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >Descrição</label
            >
            <textarea
              value={res.ProdutoDescricao}
              id="descricao"
              name="descricao"
              rows="4"
              class="block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"
              placeholder="Escreva uma descrição"
              autocomplete="off"
              required
            />
          </div>
          <div
            class="flex justify-center col-span-1 items-center pl-4 border border-gray-200 rounded dark:border-gray-700"
          >
            <img
              src={res.ProdutoImagem}
              alt={res.ProdutoNome + " imagem"}
              class="justify-self-center w-auto h-24 mr-3"
            />
          </div>
          <div class="col-span-2">
            <label
              class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              for="file_input">Enviar imagem</label
            >
            <input
              class="block w-full text-sm text-gray-900 border border-gray-300 rounded-lg cursor-pointer bg-gray-50 dark:text-gray-400 focus:outline-none dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400"
              aria-describedby="file_input_help"
              id="file_input"
              name="imagem"
              type="file"
            />
            <p
              class="mt-1 text-sm text-gray-500 dark:text-gray-300"
              id="file_input_help"
            >
              PNG, JPG or WEBP (MAX. 800x400px).
            </p>
          </div>
          <div
            class="flex items-center pl-4 border border-gray-200 rounded dark:border-gray-700"
          >
            <input
              checked={res.ProdutoAtivo}
              id="ativo"
              type="checkbox"
              name="ativo"
              class="w-4 h-4 text-red-600 bg-gray-100 border-gray-300 rounded focus:ring-red-500 dark:focus:ring-red-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
            />
            <label
              for="ativo"
              class="w-full py-4 ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
              >Ativo</label
            >
          </div>
        </div>
        <div class="flex items-center space-x-4">
          <button
            type="submit"
            class="text-white bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800"
          >
            Atualizar Produto
          </button>
          <button
            type="button"
            class="text-red-600 inline-flex items-center hover:text-white border border-red-600 hover:bg-red-600 focus:ring-4 focus:outline-none focus:ring-red-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:border-red-500 dark:text-red-500 dark:hover:text-white dark:hover:bg-red-600 dark:focus:ring-red-900"
          >
            <svg
              class="mr-1 -ml-1 w-5 h-5"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
              ><path
                fill-rule="evenodd"
                d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                clip-rule="evenodd"
              /></svg
            >
            Deletar
          </button>
        </div>
      {:else}
        <h1>Erro ao carregar produto</h1>
      {/if}
    {/await}
  </form>
</Modal>
