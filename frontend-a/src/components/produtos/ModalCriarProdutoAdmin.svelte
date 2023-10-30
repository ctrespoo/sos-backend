<script lang="ts">
  import type { Categoria } from "../../lib/types/produtos";
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
    console.log("Erro ao carregar categorias");
    toast.error("Erro ao carregar categorias");
  }

  function pegarFormulario(
    e: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }
  ) {
    const data = {
      nome: e.currentTarget.nome.value as string,
      descricao: e.currentTarget.descricao.value as string,
      peso: Number(e.currentTarget.peso.value),
      preco: Number(e.currentTarget.preco.value),
      categoria: selectedCategory,
      unidadeMedida: e.currentTarget.unidadeMedida.value as string,
      quantidadePacote: Number(e.currentTarget.qtPacote.value),
      ativo: e.currentTarget.ativo.value === "on",
    };
    const imagem = e.currentTarget.imagem.files[0] as File;

    const sendFormData = new FormData();
    sendFormData.append("imagem", imagem);
    sendFormData.append("produto", JSON.stringify(data));

    const token = localStorage.getItem("token");

    fetch(url.produtos, {
      method: "POST",
      mode: "cors",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: sendFormData,
    })
      .then(async (res) => {
        if (res.status === 201) {
          console.log(await res.text());
          toast.success("Produto criado com sucesso");
          formModal = false;
          return;
        }
        if (res.status === 200) {
          console.log(await res.text());
          toast.success("Produto criado com sucesso");
          formModal = false;
          return;
        }
        toast.error(await res.text());
      })
      .catch(async (err) => {
        toast.error("Erro ao criar produto");
        console.log(err);
      });
  }
</script>

<Modal
  bind:open={formModal}
  title="Criar Produto"
  size="md"
  autoclose={false}
  class="w-full shadow-lg shadow-slate-950"
>
  <form
    id="criarProduto"
    class="flex flex-col space-y-6"
    on:submit|preventDefault={pegarFormulario}
  >
    <!-- Modal body -->

    <div class="grid gap-4 mb-4 sm:grid-cols-3">
      <div class="col-span-3">
        <label
          for="nome"
          class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
          >Name</label
        >
        <input
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
          for="unidadeMedida"
          class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
          >Unidade de medida</label
        >
        <select
          id="unidadeMedida"
          required
          name="unidadeMedida"
          class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"
        >
          <option selected value="G">G</option>
          <option value="KG">KG</option>
        </select>
      </div>
      <div>
        <label
          for="qtPacote"
          class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
          >Quantidade por pacote</label
        >
        <input
          value={1}
          type="number"
          name="qtPacote"
          id="qtPacote"
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
          value={1}
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
          type="number"
          name="preco"
          id="preco"
          class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-primary-500 dark:focus:border-primary-500"
          placeholder="1.99"
          required
          autocomplete="off"
        />
      </div>
      <div class="col-span-2">
        <label
          class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
          for="imagem">Enviar imagem</label
        >
        <input
          class="block w-full text-sm text-gray-900 border border-gray-300 rounded-lg cursor-pointer bg-gray-50 dark:text-gray-400 focus:outline-none dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400"
          aria-describedby="imagem_help"
          id="imagem"
          name="imagem"
          type="file"
          accept="image/*"
          required
        />
        <p
          class="mt-1 text-sm text-gray-500 dark:text-gray-300"
          id="imagem_help"
        >
          PNG, JPEG (MAX/RECOMENDADO. 512x512px).
        </p>
      </div>
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
        class="flex items-center pl-4 border border-gray-200 rounded dark:border-gray-700"
      >
        <input
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
        class="text-white inline-flex items-center bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800"
      >
        <svg
          class="mr-1 -ml-1 w-6 h-6"
          fill="currentColor"
          viewBox="0 0 20 20"
          xmlns="http://www.w3.org/2000/svg"
          ><path
            fill-rule="evenodd"
            d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z"
            clip-rule="evenodd"
          /></svg
        >
        Adicionar Produto
      </button>
    </div>
  </form>
</Modal>
