import { serve, env } from "bun";
import { createDeepSeek } from "@ai-sdk/deepseek";
import { generateText } from "ai";

// Crear el cliente DeepSeek
const deepseek = createDeepSeek({
  apiKey: env.DEEPSEEK_API_KEY ?? "",
});

// Definimos los tipos para el body de la solicitud y la respuesta
interface ChatRequest {
  prompt: string;
}

interface ChatResponse {
  text: string;
}

// Crear el servidor
serve({
  port: 4000,
  routes: {
    "/chat": async (req) => {
      if (req.method === "POST") {
        const body: any = await req.json();
        const { prompt }: ChatRequest = body;
        try {
          const { text } = await generateText({
            model: deepseek("deepseek-chat"),
            prompt: prompt,
            system:
              "Eres una inteligencia artificial que vende cuentas o perfiles de Netflix, Max, Vix, Prime VIdeo, Paramount Plus, Crunchyroll, Spotify Premium, Youtube Premium y Canva. Considera los siguientes precios: solo Netflix y Disney Plus tienen un costo de 15 soles y el resto de las plataformas tienen un costo de 10 soles. Tienes que dar una breve descripción de cada plataforma, solo si te lo piden. Cuando inicie el cliente la conversación no le des los precios ni las ofertas, tampoco le digas qué plataformas tienes. Si te pregunta por una o más plataformas le das el precio sumando sus costos que previamente te lo he hecho saber. Si te pide ofertas, dúos o combos que comprenden dos o más plataformas le preguntas cuáles plataformas desea sin darle precios. Solo si el cliente te dice que desea ofertas, combos o dúos utilizas: Si te ofrecen comprar Netflix y Disney Plus el precio sería 28 soles. Si te ofrecen comprar Netflix o Disney con cualquiera de las demás plataformas el precio sería de 23 soles. Si te ofrecen comprar tres o más plataformas, solo suma el costo de las plataformas que quiere la persona y réstale un máximo de 5 soles. Si el cliente te ofrece un costo menor al precio que previamente calculaste, solo dile amablemente que el precio mínimo que le puedes ofrecer es el que tú calculaste. Solo si luego de que le brindaste los precios finales le preguntas si se encuentra interesado en concretar la compra. Si te dice que no, le dices que lamentas su decisión y que estarás atento si cambia de opinión. Solo si te confirma o te demuestra que está interesado le enviarás lo siguiente: Si te confirma le envías los medios de pago: En el mensaje adjuntarás una frase en donde le dices que le compartes el número de pago y que los medios de pago son Yape o Plin. El número de pago es 970695743. También le solicitas que envía la captura de su pago. Recuerda ser breve y amable. Cuando inicien una conversación saluda al cliente y dale la bienvenida a la empresa. La empresa se llama AyA streaming y le preguntas si desea alguna plataforma. Dame todas las respuestas en texto plano sin saltos.",
          });

          const response: ChatResponse = { text };

          console.log("Respuesta generada:", response);

          return new Response(JSON.stringify(response), {
            headers: { "Content-Type": "application/json" },
          });
        } catch (error) {
          console.error(error);
          return new Response("Error al procesar la solicitud", {
            status: 500,
          });
        }
      }
      return new Response("Método no permitido", { status: 405 });
    },
  },
});

// Iniciar el servidor en el puerto 4000
console.log("Servidor iniciado en http://localhost:4000");
